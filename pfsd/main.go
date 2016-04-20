package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/encryption"
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/pfsd/dnetclient"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/intercom"
	"github.com/cpssd/paranoid/pfsd/keyman"
	"github.com/cpssd/paranoid/pfsd/pfi"
	"github.com/cpssd/paranoid/pfsd/pnetclient"
	"github.com/cpssd/paranoid/pfsd/pnetserver"
	"github.com/cpssd/paranoid/pfsd/upnp"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	rpb "github.com/cpssd/paranoid/proto/raft"
	"github.com/cpssd/paranoid/raft"
	"github.com/cpssd/paranoid/raft/raftlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	GenerationJoinTimeout time.Duration = time.Minute * 3
	JoinSendKeysInterval  time.Duration = time.Second
)

var (
	srv *grpc.Server
	log *logger.ParanoidLogger

	certFile   = flag.String("cert", "", "TLS certificate file - if empty connection will be unencrypted")
	keyFile    = flag.String("key", "", "TLS key file - if empty connection will be unencrypted")
	skipVerify = flag.Bool("skip_verification", false,
		"skip verification of TLS certificate chain and hostname - not recommended unless using self-signed certs")
	verbose = flag.Bool("v", false, "Use verbose logging")
)

type keySentResponse struct {
	err  error
	uuid string
}

func startKeyStateMachine() {
	_, err := os.Stat(path.Join(globals.ParanoidDir, "meta", keyman.KSM_FILE_NAME))
	if err == nil {
		var err error
		keyman.StateMachine, err = keyman.NewKSMFromPFSDir(globals.ParanoidDir)
		if err != nil {
			log.Fatal("Unable to start key state machine:", err)
		}
	} else if os.IsNotExist(err) {
		keyman.StateMachine = keyman.NewKSM(globals.ParanoidDir)
	} else {
		log.Fatal("Error stating key state machine file")
	}
}

func sendKeyPiece(uuid string, piece *keyman.KeyPiece, responseChan chan keySentResponse) {
	err := pnetclient.SendKeyPiece(uuid, piece)
	responseChan <- keySentResponse{
		err:  err,
		uuid: uuid,
	}
}

func startRPCServer(lis *net.Listener, password string) {
	var opts []grpc.ServerOption
	if globals.TLSEnabled {
		log.Info("Starting ParanoidNetwork server with TLS.")
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatal("Failed to generate TLS credentials:", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	} else {
		log.Info("Starting ParanoidNetwork server without TLS.")
	}
	srv = grpc.NewServer(opts...)

	pb.RegisterParanoidNetworkServer(srv, &pnetserver.ParanoidServer{})
	nodeDetails := raft.Node{
		IP:         globals.ThisNode.IP,
		Port:       globals.ThisNode.Port,
		CommonName: globals.ThisNode.CommonName,
		NodeID:     globals.ThisNode.UUID,
	}

	startKeyStateMachine()

	if globals.Encrypted && globals.KeyGenerated {
		log.Info("Attempting to unlock")
		Unlock()
	}

	//First node to join a given cluster
	if len(globals.Nodes.GetAll()) == 0 {
		log.Info("Performing first node setup")
		globals.RaftNetworkServer = raft.NewRaftNetworkServer(
			nodeDetails,
			globals.ParanoidDir,
			path.Join(globals.ParanoidDir, "meta", "raft"),
			&raft.StartConfiguration{
				Peers: []raft.Node{},
			},
			globals.TLSEnabled,
			globals.TLSSkipVerify,
			globals.Encrypted,
		)
		timeout := time.After(GenerationJoinTimeout)
	initalGenerationLoop:
		for {
			select {
			case <-timeout:
				log.Fatal("Unable to create inital generation")
			default:
				_, _, err := globals.RaftNetworkServer.RequestNewGeneration(globals.ThisNode.UUID)
				if err == nil {
					log.Info("Successfuly created inital generation")
					break initalGenerationLoop
				}
				log.Error("Unable to create inital generation:", err)
			}
		}
		if globals.Encrypted {
			globals.KeyGenerated = true
			saveFileSystemAttributes(&globals.FileSystemAttributes{
				Encrypted:    globals.Encrypted,
				KeyGenerated: globals.KeyGenerated,
				NetworkOff:   globals.NetworkOff,
			})
		}
	} else {
		globals.RaftNetworkServer = raft.NewRaftNetworkServer(
			nodeDetails,
			globals.ParanoidDir,
			path.Join(globals.ParanoidDir, "meta", "raft"),
			nil,
			globals.TLSEnabled,
			globals.TLSSkipVerify,
			globals.Encrypted,
		)
	}

	rpb.RegisterRaftNetworkServer(srv, globals.RaftNetworkServer)

	globals.Wait.Add(1)
	go func() {
		defer globals.Wait.Done()
		err := srv.Serve(*lis)
		log.Info("Paranoid network server stopped")
		if err != nil && globals.ShuttingDown == false {
			log.Fatal("Server stopped because of an error:", err)
		}
	}()

	if globals.Encrypted && !globals.KeyGenerated {
		timeout := time.After(GenerationJoinTimeout)
	generationCreateLoop:
		for {
			select {
			case <-timeout:
				log.Fatal("Unable to join cluster before timeout")
			default:
				generation, peers, err := pnetclient.NewGeneration(password)
				if err != nil {
					log.Error("Unable to start new generation:", err)
				}

				keyPiecesN := int64(len(peers) + 1)
				keyPieces, err := keyman.GeneratePieces(globals.EncryptionKey, keyPiecesN, keyPiecesN/2+1)
				if err != nil {
					log.Fatal("Unable to split keys:", err)
				}

				err = globals.HeldKeyPieces.AddPiece(generation, globals.ThisNode.UUID, keyPieces[0])
				if err != nil {
					log.Fatal("Unable to store my key piece")
				}
				keyPieces = keyPieces[1:]

				sendKeysTimer := time.NewTimer(0)
				sendKeysResponse := make(chan keySentResponse, len(peers))
				var sendKeyPieceWait sync.WaitGroup

			sendKeysLoop:
				for {
					select {
					case <-timeout:
						log.Fatal("Unable to join cluster before timeout")
					case <-sendKeysTimer.C:
						for i := 0; i < len(peers); i++ {
							sendKeyPieceWait.Add(1)
							go func() {
								defer sendKeyPieceWait.Done()
								sendKeyPiece(peers[i], keyPieces[i], sendKeysResponse)
							}()
						}
						sendKeysTimer.Reset(JoinSendKeysInterval)
					case keySendInfo := <-sendKeysResponse:
						if keySendInfo.err != nil {
							if keySendInfo.err == keyman.ErrGenerationDeprecated {
								log.Error("Attempting to replicate keys for deprecated generation")
								break sendKeysLoop
							}
						} else {
							for i := 0; i < len(peers); i++ {
								if peers[i] == keySendInfo.uuid {
									peers = append(peers[:i], peers[i+1:]...)
									keyPieces = append(keyPieces[:i], keyPieces[i+1:]...)

									log.Info("Attempting to join raft cluster")
									err := pnetclient.JoinCluster(password)
									if err != nil {
										log.Error("Unable to join a raft cluster:", err)
									} else {
										log.Info("Sucessfully joined raft cluster")
										globals.Wait.Add(1)
										go func() {
											defer globals.Wait.Done()
											done := make(chan bool, 1)
											go func() {
												sendKeyPieceWait.Wait()
												done <- true
											}()
											for {
												select {
												case <-sendKeysResponse:
												case <-done:
													return
												}
											}
										}()
										break generationCreateLoop
									}
								}
							}
						}
					}
				}
			}
		}

		globals.KeyGenerated = true
		saveFileSystemAttributes(&globals.FileSystemAttributes{
			Encrypted:    globals.Encrypted,
			KeyGenerated: globals.KeyGenerated,
			NetworkOff:   globals.NetworkOff,
		})
	} else if globals.RaftNetworkServer.State.Configuration.HasConfiguration() == false {
		log.Info("Attempting to join raft cluster")
		err := dnetclient.JoinCluster(password)
		if err != nil {
			log.Fatal("Unable to join a raft cluster")
		}
	}
}

func setupLogging() {
	logDir := path.Join(globals.ParanoidDir, "meta", "logs")

	log = logger.New("main", "pfsd", logDir)
	dnetclient.Log = logger.New("dnetclient", "pfsd", logDir)
	pnetclient.Log = logger.New("pnetclient", "pfsd", logDir)
	pnetserver.Log = logger.New("pnetserver", "pfsd", logDir)
	upnp.Log = logger.New("upnp", "pfsd", logDir)
	keyman.Log = logger.New("keyman", "pfsd", logDir)
	raft.Log = logger.New("raft", "pfsd", logDir)
	raftlog.Log = logger.New("raftlog", "pfsd", logDir)
	commands.Log = logger.New("libpfs", "pfsd", logDir)
	intercom.Log = logger.New("intercom", "pfsd", logDir)
	globals.Log = logger.New("globals", "pfsd", logDir)

	log.SetOutput(logger.STDERR | logger.LOGFILE)
	dnetclient.Log.SetOutput(logger.STDERR | logger.LOGFILE)
	pnetclient.Log.SetOutput(logger.STDERR | logger.LOGFILE)
	pnetserver.Log.SetOutput(logger.STDERR | logger.LOGFILE)
	upnp.Log.SetOutput(logger.STDERR | logger.LOGFILE)
	keyman.Log.SetOutput(logger.STDERR | logger.LOGFILE)
	raft.Log.SetOutput(logger.STDERR | logger.LOGFILE)
	raftlog.Log.SetOutput(logger.STDERR | logger.LOGFILE)
	commands.Log.SetOutput(logger.STDERR | logger.LOGFILE)
	intercom.Log.SetOutput(logger.STDERR | logger.LOGFILE)
	globals.Log.SetOutput(logger.STDERR | logger.LOGFILE)

	if *verbose {
		commands.Log.SetLogLevel(logger.VERBOSE)
	}
}

func getFileSystemAttributes() {
	attributesJson, err := ioutil.ReadFile(path.Join(globals.ParanoidDir, "meta", "attributes"))
	if err != nil {
		log.Fatal("unable to read file system attributes:", err)
	}

	attributes := &globals.FileSystemAttributes{}
	err = json.Unmarshal(attributesJson, attributes)
	if err != nil {
		log.Fatal("unable to read file system attributes:", err)
	}

	globals.Encrypted = attributes.Encrypted
	globals.NetworkOff = attributes.NetworkOff
	encryption.Encrypted = attributes.Encrypted

	if attributes.Encrypted {
		if !attributes.KeyGenerated {
			//If a key has not yet been generated for this file system, one must be generated
			globals.EncryptionKey, err = keyman.GenerateKey(32)
			if err != nil {
				log.Fatal("unable to generate encryption key:", err)
			}

			cipherB, err := encryption.GenerateAESCipherBlock(globals.EncryptionKey.GetBytes())
			if err != nil {
				log.Fatal("unable to generate cipher block:", err)
			}
			encryption.SetCipher(cipherB)

			if attributes.NetworkOff {
				//If networking is turned off, save the key to a file
				attributes.KeyGenerated = true
				attributes.EncryptionKey = *globals.EncryptionKey
			}
		} else if attributes.NetworkOff {
			//If networking is off, load the key from the file
			globals.EncryptionKey = &attributes.EncryptionKey
			cipherB, err := encryption.GenerateAESCipherBlock(globals.EncryptionKey.GetBytes())
			if err != nil {
				log.Fatal("unable to generate cipher block:", err)
			}
			encryption.SetCipher(cipherB)
		}
	}

	globals.KeyGenerated = attributes.KeyGenerated
	saveFileSystemAttributes(attributes)
}

func saveFileSystemAttributes(attributes *globals.FileSystemAttributes) {
	attributesJson, err := json.Marshal(attributes)
	if err != nil {
		log.Fatal("unable to save new file system attributes to file:", err)
	}

	newAttributesFile := path.Join(globals.ParanoidDir, "meta", "attributes-new")
	err = ioutil.WriteFile(newAttributesFile, attributesJson, 0600)
	if err != nil {
		log.Fatal("unable to save new file system attributes to file:", err)
	}

	err = os.Rename(newAttributesFile, path.Join(globals.ParanoidDir, "meta", "attributes"))
	if err != nil {
		log.Fatal("unable to save new file system attributes to file:", err)
	}
}

func main() {
	flag.Parse()

	if len(flag.Args()) < 6 {
		fmt.Print("Usage:\n\tpfsd <paranoid_directory> <mount_point> <Discovery Server> <Discovery Port> <Discovery Pool>, <Discovery Pool Password>\n")
		os.Exit(1)
	}

	paranoidDirAbs, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		fmt.Println("FATAL: Could not get absolute paranoid dir:", err)
		os.Exit(1)
	}

	mountPointAbs, err := filepath.Abs(flag.Arg(1))
	if err != nil {
		fmt.Println("FATAL: Could not get absolute mount point:", err)
		os.Exit(1)
	}

	globals.ParanoidDir = paranoidDirAbs
	globals.MountPoint = mountPointAbs
	setupLogging()

	getFileSystemAttributes()

	globals.TLSSkipVerify = *skipVerify
	if *certFile != "" && *keyFile != "" {
		globals.TLSEnabled = true
		if !globals.TLSSkipVerify {
			cn, err := getCommonNameFromCert(*certFile)
			if err != nil {
				log.Fatal("Could not get CN from provided TLS cert:", err)
			}
			globals.ThisNode.CommonName = cn
		}
	} else {
		globals.TLSEnabled = false
	}

	if !globals.NetworkOff {
		discoveryPort, err := strconv.Atoi(flag.Arg(3))
		if err != nil || discoveryPort < 1 || discoveryPort > 65535 {
			log.Fatal("Discovery port must be a number between 1 and 65535, inclusive.")
		}

		uuid, err := ioutil.ReadFile(path.Join(globals.ParanoidDir, "meta", "uuid"))
		if err != nil {
			log.Fatal("Could not get node UUID:", err)
		}
		globals.ThisNode.UUID = string(uuid)

		ip, err := upnp.GetIP()
		if err != nil {
			log.Fatal("Could not get IP:", err)
		}

		//Asking for port 0 requests a random free port from the OS.
		lis, err := net.Listen("tcp", ip+":0")
		if err != nil {
			log.Fatalf("Failed to start listening : %v.\n", err)
		}
		splits := strings.Split(lis.Addr().String(), ":")
		port := splits[len(splits)-1]
		portInt, err := strconv.Atoi(port)
		if err != nil {
			log.Fatal("Could not parse port", splits[len(splits)-1], " Error :", err)
		}
		globals.ThisNode.Port = port

		//Try and set up uPnP. Otherwise use internal IP.
		globals.UPnPEnabled = false
		err = upnp.DiscoverDevices()
		if err == nil {
			log.Info("UPnP devices available")
			externalPort, err := upnp.AddPortMapping(ip, portInt)
			if err == nil {
				log.Info("UPnP port mapping enabled")
				port = strconv.Itoa(externalPort)
				globals.ThisNode.Port = port
				globals.UPnPEnabled = true
			}
		}

		globals.ThisNode.IP, err = upnp.GetIP()
		if err != nil {
			log.Fatal("Can't get IP. Error : ", err)
		}
		log.Info("Peer address:", globals.ThisNode.IP+":"+globals.ThisNode.Port)

		if _, err := os.Stat(globals.ParanoidDir); os.IsNotExist(err) {
			log.Fatal("Path", globals.ParanoidDir, "does not exist.")
		}
		if _, err := os.Stat(path.Join(globals.ParanoidDir, "meta")); os.IsNotExist(err) {
			log.Fatal("Path", globals.ParanoidDir, "is not valid PFS root.")
		}

		dnetclient.SetDiscovery(flag.Arg(2), flag.Arg(3))
		dnetclient.JoinDiscovery(flag.Arg(4), flag.Arg(5))
		err = globals.SetPoolPasswordHash(flag.Arg(5))
		if err != nil {
			log.Fatal("Error setting up password hash:", err)
		}
		startRPCServer(&lis, flag.Arg(5))
	}
	createPid("pfsd")
	pfi.StartPfi(*verbose)

	intercom.RunServer(path.Join(globals.ParanoidDir, "meta"))

	HandleSignals()
}

func createPid(processName string) {
	processID := os.Getpid()
	pid := []byte(strconv.Itoa(processID))
	err := ioutil.WriteFile(path.Join(globals.ParanoidDir, "meta", processName+".pid"), pid, 0600)
	if err != nil {
		log.Fatal("Failed to create PID file", err)
	}
}
