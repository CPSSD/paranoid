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
func RmdirCommand(paranoidDirectory, dirPath string) (returnCode returncodes.Code, returnError error) {
	Log.Info("rmdir command called")
	err := GetFileSystemLock(paranoidDirectory, ExclusiveLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	defer func() {
		err := UnLockFileSystem(paranoidDirectory)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
		}
	}()

	dirToDelete := GetParanoidPath(paranoidDirectory, dirPath)
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
		return returncodes.EUNEXPECTED, fmt.Errorf("error reading paranoidDirectory: %s", err)
	}

	if len(files) > 1 {
		return returncodes.ENOTEMPTY, errors.New(dirPath + " is not empty")
	}

	infoFileToDelete := path.Join(dirToDelete, "info")
	inodeBytes, code, err := GetFileInode(dirToDelete)
	if code != returncodes.OK {
		return code, err
	}

	inodeString := string(inodeBytes)

	code, err = canAccessFile(paranoidDirectory, inodeString, getAccessMode(syscall.O_WRONLY))
	if err != nil {
		return code, fmt.Errorf("unable to access %s: %s", dirPath, err)
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
