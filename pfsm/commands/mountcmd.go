package commands

import (
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

//MountCommand is used to notify a pfs directory it has been mounted.
//Stores the ip given as args[1] and the port given as args[2] in files in the meta directory.
func MountCommand(args []string) {
	Log.Info("mount command called")
	if len(args) < 4 {
		Log.Fatal("Not enough arguments!")
	}
	directory := args[0]
	Log.Verbose("mount : given directory = " + directory)

	err := ioutil.WriteFile(path.Join(directory, "meta", "ip"), []byte(args[1]), 0600)
	if err != nil {
		Log.Fatal("error writing ip:", err)
	}

	err = ioutil.WriteFile(path.Join(directory, "meta", "port"), []byte(args[2]), 0600)
	if err != nil {
		Log.Fatal("error writing port", err)
	}

	mountPoint, err := filepath.Abs(args[3])
	if err != nil {
		Log.Fatal("error getting absolute path of mountpoint:", err)
	}

	err = ioutil.WriteFile(path.Join(directory, "meta", "mountpoint"), []byte(mountPoint), 0600)
	if err != nil {
		Log.Fatal("error writing mountpoint:", err)
	}

	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
