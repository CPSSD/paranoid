package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io/ioutil"
	"os"
	"path"
)

// SymlinkCommand creates a symbolic link
func SymlinkCommand(directory, existingFile, targetFile string, sendOverNetwork bool) (returnCode int, returnError error) {
	Log.Info("symlink command called")

	targetFilePath := getParanoidPath(directory, targetFile)

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

	targetFileType, err := getFileType(directory, targetFilePath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	if targetFileType != typeENOENT {
		return returncodes.EEXIST, errors.New(targetFile + " already exists")
	}

	uuidBytes, err := generateNewInode()
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	uuidString := string(uuidBytes)
	Log.Verbose("symlink: uuid", uuidString)

	err = ioutil.WriteFile(targetFilePath, uuidBytes, 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing file:", err)
	}

	err = os.Symlink(os.DevNull, path.Join(directory, "contents", uuidString))
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error creating symlinks:", err)
	}

	nodeData := &inode{
		Inode: uuidString,
		Count: 1,
		Link:  existingFile,
	}

	jsonData, err := json.Marshal(nodeData)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error marshalling json:", err)
	}

	err = ioutil.WriteFile(path.Join(directory, "inodes", uuidString), jsonData, 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing inodes file:", err)
	}

	if sendOverNetwork {
		//Handle this later at mega binary refactor stage
		//sendToServer(directory, "symlink", args[1:], nil)
	}
	return returncodes.OK, nil
}
