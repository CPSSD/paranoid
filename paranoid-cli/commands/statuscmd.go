package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/cpssd/paranoid/pfsd/intercom"
	"io/ioutil"
	"net/rpc"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

// Status displays statistics for the specified PFSD instances.
func Status(c *cli.Context) {
	args := c.Args()
	usr, err := user.Current()
	if err != nil {
		Log.Error("Could not get user information:", err)
		fmt.Println("Unable to get information on current user:", err)
		os.Exit(1)
	}

	// By default, print the status of each running instance.
	if !args.Present() {
		dirs, err := ioutil.ReadDir(path.Join(usr.HomeDir, ".pfs", "filesystems"))
		if err != nil {
			Log.Error("Could not get list of paranoid file systems:", err)
			fmt.Printf("Unable to get list of paranoid file systems. Does %s exist?", path.Join(usr.HomeDir, ".pfs"))
			os.Exit(1)
		}
		for _, dir := range dirs {
			dirPath := path.Join(usr.HomeDir, ".pfs", "filesystems", dir.Name())
			if _, err := os.Stat(path.Join(dirPath, "meta", "pfsd.pid")); err == nil {
				getStatus(dirPath)
			}
		}
	} else {
		for _, dir := range args {
			getStatus(path.Join(usr.HomeDir, ".pfs", "filesystems", dir))
		}
	}
}

func getStatus(pfsDir string) {
	// We check this on the off chance they haven't initialised a single PFS yet.
	if _, err := os.Stat(pfsDir); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s does not exist. Please call 'paranoid-cli init' before running this command.", pfsDir)
			Log.Fatal("PFS directory does not exist.")
		} else {
			fmt.Printf("Could not stat %s. Error returned: %s.", pfsDir, err)
			Log.Fatal("Could not stat PFS directory:", err)
		}
	}

	socketPath := path.Join(pfsDir, "meta", "intercom.sock")
	logPath := path.Join(pfsDir, "meta", "logs", "pfsd.log")
	var resp intercom.StatusResponse
	client, err := rpc.Dial("unix", socketPath)
	if err != nil {
		fmt.Printf("Could not connect to PFSD %s. Is it running? See %s for more information.\n", filepath.Base(pfsDir), logPath)
		Log.Warn("Could not connect to PFSD %s at %s: %s", filepath.Base(pfsDir), socketPath, err)
		return
	}
	err = client.Call("IntercomServer.Status", new(intercom.EmptyMessage), &resp)
	if err != nil {
		fmt.Printf("Error getting status for %s. See %s for more information.\n", filepath.Base(pfsDir), logPath)
		Log.Warn("PFSD at %s returned error: %s", filepath.Base(pfsDir), err)
		return
	}
	printStatusInfo(filepath.Base(pfsDir), resp)
}

func printStatusInfo(pfsName string, info intercom.StatusResponse) {
	fmt.Printf("\nFilesystem Name:\t%s\n", pfsName)
	fmt.Printf("Uptime:\t\t\t%s\n", info.Uptime.String())
	fmt.Printf("Raft Status:\t\t%s\n", info.Status)
	fmt.Printf("TLS Enabled:\t\t%t\n", info.TLSActive)
	fmt.Printf("Port:\t\t\t%d\n", info.Port)
}
