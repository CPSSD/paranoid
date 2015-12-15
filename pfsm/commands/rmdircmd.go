package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"syscall"
)

// RmdirCommand removes a directory
func RmdirCommand(directory, dirName string, sendOverNetwork bool) (returnCode int, returnError error) {
	Log.Info("rmdir command called")

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

	dirToDelete := getParanoidPath(directory, dirName)
	dirType, err := getFileType(directory, dirToDelete)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	if dirType == typeENOENT {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return returncodes.ENOENT, errors.New(dirName + " does not exist")
	} else if dirType != typeDir {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOTDIR))
		return returncodes.ENOTDIR, errors.New(dirName + " is not a directory")
	}

	files, err := ioutil.ReadDir(dirToDelete)
	if err != nil {
		if os.IsPermission(err) {
			return returncodes.EACCES, errors.New("could not access " + dirName)
		}
		return returncodes.EUNEXPECTED, fmt.Errorf("error reading directory:", err)
	}

	if len(files) > 1 {
		return returncodes.ENOTEMPTY, errors.New(dirName + " is not empty")
	}

	infoFileToDelete := path.Join(dirToDelete, "info")
	inodeBytes, code, err := getFileInode(dirToDelete)
	if code != returncodes.OK {
		return code, err
	}

	inodeString := string(inodeBytes)

	err = syscall.Access(path.Join(directory, "contents", inodeString), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		return returncodes.EACCES, errors.New("could not access " + dirName)
	}

	inodeFileToDelete := path.Join(directory, "inodes", inodeString)
	contentsFileToDelete := path.Join(directory, "contents", inodeString)

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

	if sendOverNetwork {
		//will be sorted later at mega binary
		//sendToServer(directory, "rmdir", args[1:], nil)
	}
	return returncodes.OK, nil
}
