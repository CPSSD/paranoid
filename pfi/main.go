package pfi

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/pfi/filesystem"
	"github.com/cpssd/paranoid/pfi/util"
	"log"
	"os"
	"path/filepath"

	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

func StartPfi(pfsDir, mountPoint string, logOutput, sendOverNetwork bool) {
	// Create a logger
	var err error
	util.Log, err = logger.New("pfi", "pfi", os.DevNull)
	if err != nil {
		log.Fatalln("Error setting up logger:", err)
	}

	util.LogOutput = logOutput
	util.SendOverNetwork = sendOverNetwork
	if logOutput {
		util.Log.SetLogLevel(logger.VERBOSE)
	}

	commands.Log, err = logger.New("libpfs", "libpfs", os.DevNull)
	if err != nil {
		util.Log.Fatal("Error setting up logger:", err)
	}

	if logOutput {
		commands.Log.SetLogLevel(logger.VERBOSE)
	}

	util.PfsDirectory, err = filepath.Abs(pfsDir)
	if err != nil {
		util.Log.Fatal("Error getting pfsdirectory absoulte path :", err)
	}
	util.MountPoint, err = filepath.Abs(mountPoint)
	if err != nil {
		util.Log.Fatal("Error getting mountpoint absoulte path :", err)
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
