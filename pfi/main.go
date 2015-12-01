package main

import (
	"flag"
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/pfi/filesystem"
	"github.com/cpssd/paranoid/pfi/pfsminterface"
	"github.com/cpssd/paranoid/pfi/util"
	"log"
	"path/filepath"

	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

func main() {
	// parsing flags and args
	logOutput := flag.Bool("v", false, "Log operations in standard output")
	markNetwork := flag.Bool("n", false, "Mark file system operations as coming from the network")
	flag.Parse()

	// Create a logger
	util.Log = logger.New("main", "pfi", "/dev/null")
	if *logOutput {
		util.Log.SetLogLevel(logger.VERBOSE)
	}

	if *markNetwork {
		pfsminterface.OriginFlag = "-n"
	} else {
		pfsminterface.OriginFlag = "-f"
	}
	noFlagArgs := flag.Args()

	if len(noFlagArgs) < 2 {
		log.Fatalln("\nUsage:\npfi [flags] <PfsInitPoint> <MountPoint>")
	}

	var err error
	util.PfsDirectory, err = filepath.Abs(noFlagArgs[0])
	if err != nil {
		log.Fatalln(err)
	}
	util.MountPoint, err = filepath.Abs(noFlagArgs[1])
	if err != nil {
		log.Fatalln(err)
	}

	// setting up with fuse
	opts := pathfs.PathNodeFsOptions{}
	opts.ClientInodes = true
	nfs := pathfs.NewPathNodeFs(&filesystem.ParanoidFileSystem{
		FileSystem: pathfs.NewDefaultFileSystem(),
	}, &opts)
	server, _, err := nodefs.MountRoot(util.MountPoint, nfs.Root(), nil)
	if err != nil {
		util.Log.Fatalf("Mount fail: %v\n", err)
	}
	server.Serve()
}
