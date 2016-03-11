package commands

import (
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"io/ioutil"
	"path"
	"path/filepath"
)

//MountCommand is used to notify a pfs paranoidDirectory it has been mounted.
func MountCommand(paranoidDirectory, ip, port, mountPoint string) (returnCode int, returnError error) {
	Log.Verbose("mount command called")
	Log.Verbose("mount : given paranoidDirectory = " + paranoidDirectory)

	err := ioutil.WriteFile(path.Join(paranoidDirectory, "meta", "ip"), []byte(ip), 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing ip:", err)
	}

	err = ioutil.WriteFile(path.Join(paranoidDirectory, "meta", "port"), []byte(port), 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing port:", err)
	}

	mountPoint, err = filepath.Abs(mountPoint)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error getting absolute path of mountpoint:", err)
	}

	err = ioutil.WriteFile(path.Join(paranoidDirectory, "meta", "mountpoint"), []byte(mountPoint), 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing mountpoint:", err)
	}

	return returncodes.OK, nil
}
