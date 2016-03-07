package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"io/ioutil"
	"os"
	"path"
	"syscall"
)

// RmdirCommand removes a paranoidDirectory
func RmdirCommand(paranoidDirectory, dirPath string) (returnCode int, returnError error) {
	Log.Info("rmdir command called")

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

	dirToDelete := getParanoidPath(paranoidDirectory, dirPath)
	dirType, err := getFileType(paranoidDirectory, dirToDelete)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	if dirType == typeENOENT {
		return returncodes.ENOENT, errors.New(dirPath + " does not exist")
	} else if dirType != typeDir {
		return returncodes.ENOTDIR, errors.New(dirPath + " is not a paranoidDirectory")
	}

	files, err := ioutil.ReadDir(dirToDelete)
	if err != nil {
		if os.IsPermission(err) {
			return returncodes.EACCES, errors.New("could not access " + dirPath)
		}
		return returncodes.EUNEXPECTED, fmt.Errorf("error reading paranoidDirectory:", err)
	}

	if len(files) > 1 {
		return returncodes.ENOTEMPTY, errors.New(dirPath + " is not empty")
	}

	infoFileToDelete := path.Join(dirToDelete, "info")
	inodeBytes, code, err := getFileInode(dirToDelete)
	if code != returncodes.OK {
		return code, err
	}

	inodeString := string(inodeBytes)

	err = syscall.Access(path.Join(paranoidDirectory, "contents", inodeString), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		return returncodes.EACCES, errors.New("could not access " + dirPath)
	}

	inodeFileToDelete := path.Join(paranoidDirectory, "inodes", inodeString)
	contentsFileToDelete := path.Join(paranoidDirectory, "contents", inodeString)

	code, err = deleteFile(contentsFileToDelete)
	if code != returncodes.OK {
		return code, err
	}

	code, err = deleteFile(inodeFileToDelete)
	if code != returncodes.OK {
		return code, err
	}

	code, err = deleteFile(infoFileToDelete)
	if code != returncodes.OK {
		return code, err
	}

	code, err = deleteFile(dirToDelete)
	if code != returncodes.OK {
		return code, err
	}

	return returncodes.OK, nil
}
