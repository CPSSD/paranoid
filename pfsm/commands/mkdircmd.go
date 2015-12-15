package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"os"
	"path"
)

// MkdirCommand is called when making a directory
func MkdirCommand(directory, dirName string, mode os.FileMode, sendOverNetwork bool) (returnCode int, returnError error) {
	Log.Info("mkdir command called")

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

	dirPath := getParanoidPath(directory, dirName)
	dirInfoPath := path.Join(dirPath, "info")

	inodeBytes, err := generateNewInode()
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	inodeString := string(inodeBytes)
	inodePath := path.Join(directory, "inodes", inodeString)
	contentsPath := path.Join(directory, "contents", inodeString)

	fileType, err := getFileType(directory, dirPath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	if fileType != typeENOENT {
		return returncodes.EEXIST, errors.New(dirName + " already exists")
	}

	err = os.Mkdir(dirPath, mode)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error making directory "+dirPath+" :", err)
	}

	contentsFile, err := os.Create(contentsPath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error creating contents file:", err)
	}

	err = contentsFile.Chmod(os.FileMode(mode))
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error changing file permissions:", err)
	}

	dirInfoFile, err := os.Create(dirInfoPath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error creating info file:", err)
	}

	_, err = dirInfoFile.Write(inodeBytes)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing to info file:", err)
	}

	inodeFile, err := os.Create(inodePath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error creating inode file:", err)
	}

	nodeData := &inode{
		Inode: inodeString,
		Count: 1}
	jsonData, err := json.Marshal(nodeData)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error marshalling json:", err)
	}

	_, err = inodeFile.Write(jsonData)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing to inode file:", err)
	}

	if sendOverNetwork {
		//This will be sorted later when we get rid of IC
		//sendToServer(directory, "mkdir", args[1:], nil)
	}
	return returncodes.OK, nil
}
