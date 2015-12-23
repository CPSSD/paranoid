package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"io/ioutil"
	"os"
	"path"
)

//CreatCommand creates a new file with the name fileName in the pfs directory
func CreatCommand(directory, fileName string, perms os.FileMode, sendOverNetwork bool) (returnCode int, returnError error) {
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

	fileType, err := getFileType(directory, namepath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	if fileType != typeENOENT {
		return returncodes.EEXIST, errors.New(fileName + " already exists")
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
		return returncodes.EUNEXPECTED, fmt.Errorf("error marshalling json:", err)
	}

	err = ioutil.WriteFile(path.Join(directory, "inodes", uuidstring), jsonData, 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing inodes file:", err)
	}

	contentsFile, err := os.Create(path.Join(directory, "contents", uuidstring))
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error creating contents file:", err)
	}

	err = contentsFile.Chmod(os.FileMode(perms))
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error changing file permissions:", err)
	}

	if sendOverNetwork {
		//This will be sorted later when we get rid of IC
		//sendToServer(directory, "creat", args[1:], nil)
	}
	return returncodes.OK, nil
}
