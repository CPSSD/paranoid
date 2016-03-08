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
		dirs, err := ioutil.ReadDir(path.Join(usr.HomeDir, ".pfs"))
		if err != nil {
			Log.Error("Could not get list of paranoid file systems:", err)
			fmt.Printf("Unable to get list of paranoid file systems. Does %s exist?", path.Join(usr.HomeDir, ".pfs"))
			os.Exit(1)
		}
		for _, dir := range dirs {
			dirPath := path.Join(usr.HomeDir, ".pfs", dir.Name())
			if _, err := os.Stat(path.Join(dirPath, "meta", "pfsd.pid")); err != nil {
				getStatus(dirPath)
			}
		}
	} else {
		for _, dir := range args {
			getStatus(path.Join(usr.HomeDir, ".pfs", dir))
		}
	}
}

func getStatus(pfsDir string) {
	socketPath := path.Join(pfsDir, "meta", "intercom.sock")
	logPath := path.Join(pfsDir, "meta", "logs", "pfsd.log")
	var resp intercom.StatusResponse
	client, err := rpc.Dial("unix", socketPath)
	if err != nil {
		fmt.Printf("Could not connect to PFSD. Is it running? See %s for more information.\n", logPath)
		Log.Warn("Could not connect to PFSD at", socketPath)
		os.Exit(1)
	}
	err = client.Call("IntercomServer.Status", new(intercom.EmptyMessage), &resp)
	if err != nil {
		fmt.Printf("Error getting status for %s. See %s for more information.\n", pfsDir, logPath)
		Log.Warn("PFSD at %s returned error: %s", pfsDir, err)
		os.Exit(1)
	}
	printStatusInfo(resp)
}

func printStatusInfo(info intercom.StatusResponse) {
	fmt.Printf("Uptime:\t\t%s\n", info.Uptime.String())
	var statusString string
	switch info.Status {
	case intercom.FOLLOWER:
		statusString = "Follower"
	case intercom.CANDIDATE:
		statusString = "Candidate"
	case intercom.LEADER:
		statusString = "Leader"
	case intercom.RAFT_INACTIVE:
		statusString = "Raft Inactive"
	}
	fmt.Printf("Raft Status:\t%s\n", statusString)
	fmt.Printf("TLS Enabled:\t%t\n", info.TLSActive)
	fmt.Printf("Port:\t\t%d\n", info.Port)
}
