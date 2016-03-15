package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"io/ioutil"
	"os"
	"path"
	"syscall"
)

// SymlinkCommand creates a symbolic link
func SymlinkCommand(paranoidDirectory, existingFilePath, targetFilePath string) (returnCode int, returnError error) {
	Log.Info("symlink command called")
	targetParanoidPath := getParanoidPath(paranoidDirectory, targetFilePath)

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

	targetFilePathType, err := getFileType(paranoidDirectory, targetParanoidPath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	if targetFilePathType != typeENOENT {
		return returncodes.EEXIST, errors.New(targetFilePath + " already exists")
	}

	uuidBytes, err := generateNewInode()
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	uuidString := string(uuidBytes)
	Log.Verbose("symlink: uuid", uuidString)

	err = ioutil.WriteFile(targetParanoidPath, uuidBytes, 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing file:", err)
	}

	contentsFile := path.Join(paranoidDirectory, "contents", uuidString)
	err = os.Symlink(os.DevNull, contentsFile)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error creating symlinks:", err)
	}

	fi, err := os.Lstat(contentsFile)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error Lstating file:", err)
	}

	stat := fi.Sys().(*syscall.Stat_t)

	nodeData := &inode{
		Mode:  os.FileMode(stat.Mode),
		Inode: uuidString,
		Count: 1,
		Link:  existingFilePath,
	}

	jsonData, err := json.Marshal(nodeData)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error marshalling json:", err)
	}

	err = ioutil.WriteFile(path.Join(paranoidDirectory, "inodes", uuidString), jsonData, 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing inodes file:", err)
	}

	return returncodes.OK, nil
}
