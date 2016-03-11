package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"os"
	"path"
)

// MkdirCommand is called when making a paranoidDirectory
func MkdirCommand(paranoidDirectory, dirPath string, mode os.FileMode) (returnCode int, returnError error) {
	Log.Verbose("mkdir command called")
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

	dirParanoidPath := getParanoidPath(paranoidDirectory, dirPath)
	dirInfoPath := path.Join(dirParanoidPath, "info")

	inodeBytes, err := generateNewInode()
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	inodeString := string(inodeBytes)
	inodePath := path.Join(paranoidDirectory, "inodes", inodeString)
	contentsPath := path.Join(paranoidDirectory, "contents", inodeString)

	fileType, err := getFileType(paranoidDirectory, dirParanoidPath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	if fileType != typeENOENT {
		return returncodes.EEXIST, errors.New(dirPath + " already exists")
	}

	err = os.Mkdir(dirParanoidPath, mode)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error making paranoidDirectory "+dirParanoidPath+" :", err)
	}

	contentsFile, err := os.Create(contentsPath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error creating contents file:", err)
	}
	defer contentsFile.Close()

	err = contentsFile.Chmod(os.FileMode(mode))
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error changing file permissions:", err)
	}

	dirInfoFile, err := os.Create(dirInfoPath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error creating info file:", err)
	}
	defer dirInfoFile.Close()

	_, err = dirInfoFile.Write(inodeBytes)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing to info file:", err)
	}

	inodeFile, err := os.Create(inodePath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error creating inode file:", err)
	}
	defer inodeFile.Close()

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

	return returncodes.OK, nil
}
