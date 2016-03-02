package main

import (
	"flag"
	"github.com/cpssd/paranoid/discovery-server/dnetserver"
	"github.com/cpssd/paranoid/logger"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
	"os/user"
	"path"
	"strconv"
	"syscall"
	"time"
)

var (
	port          = flag.Int("port", 10101, "port to listen on")
	logDir        = flag.String("log-directory", "/var/log", "directory in which to create ParanoidDiscovery.log")
	renewInterval = flag.Int("renew-interval", 5*60*1000, "time after which membership expires, in ms")
	certFile      = flag.String("cert", "", "TLS certificate file - if empty connection will be unencrypted")
	keyFile       = flag.String("key", "", "TLS key file - if empty connection will be unencrypted")
	loadState     = flag.Bool("state", true, "Load the Nodes from the statefile")
)

func createRPCServer() *grpc.Server {
	var opts []grpc.ServerOption
	if *certFile != "" && *keyFile != "" {
		dnetserver.Log.Info("Starting discovery server with TLS.")
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			dnetserver.Log.Fatal("Failed to generate TLS credentials:", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	} else {
		dnetserver.Log.Info("Starting discovery server without TLS.")
	}
	return grpc.NewServer(opts...)
}

func main() {
	flag.Parse()
	dnetserver.Log = logger.New("main", "discovery-server", *logDir)
	err := dnetserver.Log.SetOutput(logger.LOGFILE | logger.STDERR)
	if err != nil {
		dnetserver.Log.Error("Failed to set logger output:", err)
	}

	analyseWorkspace(dnetserver.Log)

	renewDuration, err := time.ParseDuration(strconv.Itoa(*renewInterval) + "ms")
	if err != nil {
		dnetserver.Log.Error("Failed parsing renew interval", err)
	}

	dnetserver.RenewInterval = renewDuration

	if *port < 1 || *port > 65535 {
		dnetserver.Log.Fatal("Port must be a number between 1 and 65535, inclusive.")
	}

	dnetserver.Log.Info("Starting Paranoid Discovery Server")

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		dnetserver.Log.Fatalf("Failed to listen on port %d: %v.", *port, err)
	}
	dnetserver.Log.Info("Listening on port", *port)

	if *loadState {
		dnetserver.LoadState()
	}

	srv := createRPCServer()
	pb.RegisterDiscoveryNetworkServer(srv, &dnetserver.DiscoveryServer{})

	dnetserver.Log.Info("gRPC server created")
	srv.Serve(lis)
}

// analyseWorkspace analyses the state of the workspace directory for the server,
// if the workspace directory doesnt exist it will be recreated along with needed
// sub-directories.
func analyseWorkspace(log *logger.ParanoidLogger) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("Couldn't identify user:", err)
	}

	// checking ~/.pfs
	pfsDirPath := path.Join(usr.HomeDir, ".pfs")
	checkDir(pfsDirPath, log)

	// checking ~/.pfs/discovery_meta
	metaDirPath := path.Join(pfsDirPath, "discovery_meta")
	checkDir(metaDirPath, log)

	dnetserver.StateFilePath = path.Join(metaDirPath, "server_state.json")
}

// checkDir checks a directory and creates it if needed
func checkDir(dir string, log *logger.ParanoidLogger) {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Info("Creating: ", dir)
			err1 := os.Mkdir(dir, 0700)
			if err1 != nil {
				log.Fatal("Failed to create: ", dir, " err:", err1)
			}
		} else {
			log.Fatal("Couldn't stat:", dir, "err:", err)
		}
	} else {
		err = syscall.Access(dir, syscall.O_RDWR)
		if err != nil {
			log.Fatal("Don't have read & write access to:", dir)
		}
	}
}
