package pfi

import (
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/pfi/filesystem"
	"github.com/cpssd/paranoid/pfi/util"
	"path/filepath"

	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

func StartPfi(pfsDir, mountPoint string, logOutput, sendOverNetwork bool) {
	// Create a logger
	util.Log, _ = logger.New("pfi", "pfi", "/dev/null")
	util.LogOutput = logOutput
	util.SendOverNetwork = sendOverNetwork
	if logOutput {
		util.Log.SetLogLevel(logger.VERBOSE)
	}

	var err error
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
