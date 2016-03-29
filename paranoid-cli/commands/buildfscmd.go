package commands

import (
	"fmt"
	progress "github.com/cheggaaa/pb"
	"github.com/codegangsta/cli"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	pb "github.com/cpssd/paranoid/proto/raft"
	"github.com/cpssd/paranoid/raft"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
)

// Buildfs creates a new filesystem from a set of activity logs specified
func Buildfs(c *cli.Context) {
	args := c.Args()
	if len(args) != 2 {
		cli.ShowCommandHelp(c, "buildfs")
		os.Exit(1)
	}

	pfsName := args[0]
	logDir := args[1]

	if fileSystemExists(pfsName) {
		fmt.Println(pfsName, "already exists. please chose a different name")
		os.Exit(1)
	}

	logs, err := ioutil.ReadDir(logDir)
	if err != nil {
		fmt.Println("Couldn't read log-directory")
		log.Fatal(err)
	}

	if len(logs) < 1 {
		fmt.Println("log-directory empty")
		os.Exit(1)
	}

	doInit(pfsName, c.String("pool"), c.String("cert"),
		c.String("key"), c.Bool("unsecure"), c.Bool("unencrypted"), c.Bool("networkoff"))

	u, err := user.Current()
	if err != nil {
		fmt.Println("Couldn't get current user home directory")
		log.Fatal(err)
	}

	pfsPath := path.Join(u.HomeDir, ".pfs", "filesystems", pfsName)
	err = os.MkdirAll(path.Join(pfsPath, "meta", "raft", "raft_logs"), 0700)
	if err != nil {
		cleanupPFS(pfsPath)
		log.Fatalln(err)
	}

	bar := progress.StartNew(len(logs))
	for _, lg := range logs {
		logEntry, err := fileToProto(lg, logDir)
		if err != nil {
			cleanupPFS(pfsPath)
			fmt.Println("broken file in log directory: ", lg)
			log.Fatal(err)
		}

		if logEntry.Entry.Type == pb.Entry_StateMachineCommand {
			er := raft.PerformLibPfsCommand(pfsPath, logEntry.Entry.Command)
			if er.Code == returncodes.EUNEXPECTED {
				fmt.Println("pfsLib failed on command: ", logEntry.Entry.Command)
				log.Fatal(er.Err)
			}
		}
		bar.Increment()
	}

	bar.FinishPrint("Done.\nYou may now mount your new filesystem: " + pfsName)
}
