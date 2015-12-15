package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"os"
	"path"
	"syscall"
)

// RenameCommand is called when renaming a file
func RenameCommand(directory, oldFileName, newFileName string, sendOverNetwork bool) (returnCode int, returnError error) {
	Log.Info("rename command called")

	oldFilePath := getParanoidPath(directory, oldFileName)
	newFilePath := getParanoidPath(directory, newFileName)

	oldFileType, err := getFileType(directory, oldFilePath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error getting "+oldFileName+" type:", err)
	}

	newFileType, err := getFileType(directory, newFilePath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error getting "+newFileName+" type:", err)
	}

	err = getFileSystemLock(directory, exclusiveLock)
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

	if oldFileType == typeENOENT {
		return returncodes.ENOENT, errors.New(oldFileName + " does not exist")
	}

	if newFileType != typeENOENT {
		return returncodes.EEXIST, errors.New(newFileName + " already exists")
	}

	inodeBytes, code, err := getFileInode(oldFilePath)
	if code != returncodes.OK || err != nil {
		return code, err
	}

	err = syscall.Access(path.Join(directory, "contents", string(inodeBytes)), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		return returncodes.EACCES, errors.New("can not access " + oldFileName)
	}

	err = os.Rename(oldFilePath, newFilePath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error renaming file:", err)
	}

	if sendOverNetwork {
		//This will be sorted later when we get rid of IC
		//sendToServer(directory, "rename", args[1:], nil)
	}
	return returncodes.OK, nil
}
