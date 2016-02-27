package commands

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/pfsd/intercom"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"sync"
	"time"
)

//Talks to the necessary other programs to mount the pfs filesystem.
//If the file system doesn't exist it creates it.
func Mount(c *cli.Context) {
	args := c.Args()
	if len(args) < 2 {
		cli.ShowCommandHelp(c, "mount")
		os.Exit(1)
	}
	doMount(c, args)
}

func doMount(c *cli.Context, args []string) {
	var serverAddress string
	if serverAddress = c.String("discovery-addr"); len(serverAddress) == 0 {
		serverAddress = "paranoid.discovery.razoft.net:10101"
	}
	pfsName := args[0]
	mountPoint := args[1]

	if c.GlobalBool("networkoff") == false {
		_, err := net.DialTimeout("tcp", serverAddress, time.Duration(5*time.Second))
		if err != nil {
			Log.Fatal("Unable to reach server", err)
		}
	}

	usr, err := user.Current()
	if err != nil {
		Log.Fatal(err)
	}
	pfsDir := path.Join(usr.HomeDir, ".pfs", pfsName)

	if _, err := os.Stat(pfsDir); os.IsNotExist(err) {
		Log.Fatal("PFS directory does not exist")
	}
	if _, err := os.Stat(path.Join(pfsDir, "contents")); os.IsNotExist(err) {
		Log.Fatal("PFS directory does not include contents directory")
	}
	if _, err := os.Stat(path.Join(pfsDir, "meta")); os.IsNotExist(err) {
		Log.Fatal("PFS directory does not include meta directory")
	}
	if _, err := os.Stat(path.Join(pfsDir, "names")); os.IsNotExist(err) {
		Log.Fatal("PFS directory does not include names directory")
	}
	if _, err := os.Stat(path.Join(pfsDir, "inodes")); os.IsNotExist(err) {
		Log.Fatal("PFS directory does not include inodes directory")
	}

	if pathExists(path.Join(pfsDir, "meta/", "pfsd.pid")) {
		err = os.Remove(path.Join(pfsDir, "meta/", "pfsd.pid"))
		if err != nil {
			Log.Fatal("Could not remove old pfsd.pid")
		}
	}

	poolBytes, err := ioutil.ReadFile(path.Join(pfsDir, "meta", "pool"))
	if err != nil {
		Log.Fatal("unable to read pool information:", err)
	}
	pool := string(poolBytes)

	splitAddress := strings.Split(serverAddress, ":")
	if len(splitAddress) != 2 {
		Log.Fatal("discovery address in wrong format. Should be HOST:PORT")
	}

	returncode, err := commands.MountCommand(pfsDir, splitAddress[0], splitAddress[1], mountPoint)
	if returncode != returncodes.OK {
		Log.Fatal("Error running pfs mount command : ", err)
	}

	if !c.GlobalBool("networkoff") {
		// Check if the cert and key files are present.
		certPath := path.Join(pfsDir, "meta", "cert.pem")
		keyPath := path.Join(pfsDir, "meta", "key.pem")
		if pathExists(certPath) && pathExists(keyPath) {
			Log.Info("Starting PFSD in secure mode.")
			//TODO(terry): Add a way to check if the given cert is its own CA,
			// and skip validation based on that.
			pfsdArgs := []string{"-cert=" + certPath, "-key=" + keyPath, "-skip_verification",
				pfsDir, mountPoint, splitAddress[0], splitAddress[1], pool}
			var pfsdFlags []string
			if c.GlobalBool("verbose") {
				pfsdFlags = append(pfsdFlags, "-v")
			}
			iface := c.String("interface")
			if iface != "" {
				pfsdFlags = append(pfsdFlags, "-interface="+iface)
			}
			cmd := exec.Command("pfsd", append(pfsdFlags, pfsdArgs...)...)
			err = cmd.Start()
			if err != nil {
				Log.Fatal("Error running pfsd command :", err)
			}
		} else {
			// Start in unsecure mode
			if !c.Bool("noprompt") {
				scanner := bufio.NewScanner(os.Stdin)
				fmt.Print("Starting networking in unsecure mode. Are you sure? [y/N] ")
				scanner.Scan()
				answer := strings.ToLower(scanner.Text())
				if !strings.HasPrefix(answer, "y") {
					fmt.Println("Okay. Exiting ...")
					os.Exit(1)
				}
			}

			Log.Info("Starting PFSD in unsecure mode.")
			pfsdArgs := []string{pfsDir, mountPoint, splitAddress[0], splitAddress[1], pool}
			var pfsdFlags []string
			if c.GlobalBool("verbose") {
				pfsdFlags = append(pfsdFlags, "-v")
			}
			iface := c.String("interface")
			if iface != "" {
				pfsdFlags = append(pfsdFlags, "-interface="+iface)
			}
			cmd := exec.Command("pfsd", append(pfsdFlags, pfsdArgs...)...)
			err = cmd.Start()
			if err != nil {
				Log.Fatal("Error running pfsd :", err)
			}
		}
	} else {
		//No need to worry about security certs
		pfsdArgs := []string{pfsDir, mountPoint, splitAddress[0], splitAddress[1], pool}
		var pfsdFlags []string
		if c.GlobalBool("verbose") {
			pfsdFlags = append(pfsdFlags, "-v")
		}
		pfsdFlags = append(pfsdFlags, "--no_networking")
		cmd := exec.Command("pfsd", append(pfsdFlags, pfsdArgs...)...)
		err = cmd.Start()
		if err != nil {
			Log.Fatal("Error running pfsd :", err)
		}
	}
	// Now that we've successfully told PFSD to start, ping it until we can confirm it is up
	var ws sync.WaitGroup
	ws.Add(1)
	go func() {
		defer ws.Done()
		socketPath := path.Join(pfsDir, "meta", "intercom.sock")
		after := time.After(time.Second * 10)
		for {
			select {
			case <-after:
				fmt.Printf("PFSD failed to start: received no response from PFSD. See %s for more information.\n",
					path.Join(pfsDir, "meta", "logs", "pfsd.log"))
				return
			default:
				var resp intercom.EmptyMessage
				client, err := rpc.Dial("unix", socketPath)
				if err != nil {
					time.Sleep(time.Second)
					continue
				}
				err = client.Call("IntercomServer.ConfirmUp", new(intercom.EmptyMessage), &resp)
				if err == nil {
					return
				}
				time.Sleep(time.Second)
			}
		}
	}()
	ws.Wait()
}
