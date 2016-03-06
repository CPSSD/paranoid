package pfi

import (
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/pfsd/globals"
	"path"

	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

func StartPfi(logOutput bool) {
	defer globals.Wait.Done()
	// Create a logger
	var err error
	Log = logger.New("pfi", "pfsd", path.Join(globals.ParanoidDir, "meta", "logs"))
	Log.SetOutput(logger.STDERR | logger.LOGFILE)

	LogOutput = logOutput
	if globals.RaftNetworkServer == nil {
		SendOverNetwork = false
	} else {
		SendOverNetwork = true
	}

	if logOutput {
		Log.SetLogLevel(logger.VERBOSE)
	}

	// setting up with fuse
	opts := pathfs.PathNodeFsOptions{}
	opts.ClientInodes = true
	nfs := pathfs.NewPathNodeFs(&ParanoidFileSystem{
		FileSystem: pathfs.NewDefaultFileSystem(),
	}, &opts)
	server, _, err := nodefs.MountRoot(globals.MountPoint, nfs.Root(), nil)
	if err != nil {
		Log.Fatalf("Mount fail: %v\n", err)
	}
	go server.Serve()

	select {
	case _, ok := <-globals.Quit:
		if !ok {
			err = server.Unmount()
			if err != nil {
				Log.Fatal("Error unmounting : ", err)
			}
		}
	}
}
