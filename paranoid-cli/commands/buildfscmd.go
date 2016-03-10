package commands

import (
	progress "github.com/cheggaaa/pb"
	"github.com/codegangsta/cli"
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
		showErrAndExit(nil, pfsName, "already exists. please chose a different name")
	}

	logs, err := ioutil.ReadDir(logDir)
	if err != nil {
		showErrAndExit(err, "Couldn't read log-directory, err: ", err)
	}

	if len(logs) < 1 {
		showErrAndExit(nil, "log-directory empty")
	}

	doInit(pfsName, c.String("pool"), c.String("cert"),
		c.String("key"), c.Bool("unsecure"))

	u, err := user.Current()
	if err != nil {
		showErrAndExit(err, "Couldn't get current user err:", err)
	}

	pfsPath := path.Join(u.HomeDir, ".pfs", pfsName)

	os.MkdirAll(path.Join(pfsPath, "meta", "raft", "raft_logs"), 0777)

	bar := progress.StartNew(len(logs))
	for _, lg := range logs {
		logEntry, err := fileToProto(lg, logDir)
		if err != nil {
			cleanupPFS(pfsPath)
			showErrAndExit(err, "broken file in log-dir: ", lg)
		}

		if logEntry.Entry.Type == pb.Entry_StateMachineCommand {
			er := raft.PerformLibPfsCommand(pfsPath, logEntry.Entry.Command)
			if er.Err != nil {
				showErrAndExit(er.Err, "pfsLib failed on command: ", logEntry.Entry.Command, "With ")
			}
		}
		bar.Increment()
	}

	bar.FinishPrint("Done.\nYou may now mount your new filesystem: " + pfsName)
}

// showErrAndExit prints the error, show command help and exit.
func showErrAndExit(err error, v ...interface{}) {
	log.Println(v)
	if err == nil {
		os.Exit(1)
	} else {
		log.Fatal(err)
	}
}
