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

//CreatCommand creates a new file with the name filePath in the pfs paranoidDirectory
func CreatCommand(paranoidDirectory, filePath string, perms os.FileMode) (returnCode int, returnError error) {
	Log.Info("creat command called")
	Log.Verbose("creat : paranoidDirectory = " + paranoidDirectory)

	err := getFileSystemLock(paranoidDirectory, exclusiveLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	defer func() {
		err := unLockFileSystem(paranoidDirectory)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
		}
	}()

	namepath := getParanoidPath(paranoidDirectory, filePath)

	fileType, err := getFileType(paranoidDirectory, namepath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	if fileType != typeENOENT {
		return returncodes.EEXIST, errors.New(filePath + " already exists")
	}
	Log.Verbose("creat : creating file " + filePath)

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
		Mode:  perms,
		Inode: uuidstring,
		Count: 1}
	jsonData, err := json.Marshal(nodeData)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error marshalling json:", err)
	}

	err = ioutil.WriteFile(path.Join(paranoidDirectory, "inodes", uuidstring), jsonData, 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing inodes file:", err)
	}

	contentsFile, err := os.Create(path.Join(paranoidDirectory, "contents", uuidstring))
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error creating contents file:", err)
	}
	defer contentsFile.Close()

	return returncodes.OK, nil
}
