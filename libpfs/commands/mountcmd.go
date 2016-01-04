package commands

import (
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"io/ioutil"
	"path"
	"path/filepath"
)

//MountCommand is used to notify a pfs directory it has been mounted.
//Stores the ip given as args[1] and the port given as args[2] in files in the meta directory.
func MountCommand(directory, ip, port, mountPoint string) (returnCode int, returnError error) {
	Log.Info("mount command called")
	Log.Verbose("mount : given directory = " + directory)

	err := ioutil.WriteFile(path.Join(directory, "meta", "ip"), []byte(ip), 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing ip:", err)
	}

	err = ioutil.WriteFile(path.Join(directory, "meta", "port"), []byte(port), 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing port:", err)
	}

	mountPoint, err = filepath.Abs(mountPoint)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error getting absolute path of mountpoint:", err)
	}

	err = ioutil.WriteFile(path.Join(directory, "meta", "mountpoint"), []byte(mountPoint), 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing mountpoint:", err)
	}

	return returncodes.OK, nil
}
