package commands

import (
	"encoding/json"
	"errors"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io/ioutil"
	"os"
	"path"
)

//CreatCommand creates a new file with the name fileName in the pfs directory
func CreatCommand(directory, fileName string, perms os.FileMode) (returnCode int, returnError error) {
	Log.Info("creat command called")
	Log.Verbose("creat : directory = " + directory)

	err := getFileSystemLock(directory, exclusiveLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	defer func() {
		err := unLockFileSystem(directory)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
		}
	}()

	namepath := getParanoidPath(directory, fileName)

	if getFileType(directory, namepath) != typeENOENT {
		return returncodes.EEXIST, errors.New("file already exists")
	}
	Log.Verbose("creat : creating file " + fileName)

	uuidbytes, err := generateNewInode()
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	uuidstring := string(uuidbytes)
	Log.Verbose("creat : uuid = " + uuidstring)

	err = ioutil.WriteFile(namepath, uuidbytes, 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, errors.New("error writing name file")
	}

	nodeData := &inode{
		Inode: uuidstring,
		Count: 1}
	jsonData, err := json.Marshal(nodeData)
	if err != nil {
		return returncodes.EUNEXPECTED, errors.New("error marshalling json:", err)
	}

	err = ioutil.WriteFile(path.Join(directory, "inodes", uuidstring), jsonData, 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, errors.New("error writing inodes file:", err)
	}

	contentsFile, err := os.Create(path.Join(directory, "contents", uuidstring))
	if err != nil {
		return returncodes.EUNEXPECTED, errors.New("error creating contents file:", err)
	}

	err = contentsFile.Chmod(os.FileMode(perms))
	if err != nil {
		return returncodes.EUNEXPECTED, errors.New("error changing file permissions:", err)
	}

	if !Flags.Network {
		sendToServer(directory, "creat", args[1:], nil)
	}
	return returncodes.OK, nil
}
