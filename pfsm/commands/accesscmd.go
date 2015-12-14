package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"path"
	"syscall"
)

//AccessCommand is used by fuse to check if it has access to a given file.
func AccessCommand(directory, fileName string, mode uint32) (returnCode int, returnError error) {
	Log.Info("access command given")
	Log.Verbose("access : given directory = " + directory)

	err := getFileSystemLock(directory, sharedLock)
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

	namePath := getParanoidPath(directory, fileName)

	fileType, err := getFileType(directory, namePath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error getting "+fileName+" file type:", err)
	}

	if fileType == typeENOENT {
		return returncodes.ENOENT, errors.New(fileName + " does not exist")
	}

	inodeNameBytes, code, err := getFileInode(namePath)
	if code != returncodes.OK || err != nil {
		return code, err
	}

	inodeName := string(inodeNameBytes)

	err = syscall.Access(path.Join(directory, "contents", inodeName), mode)
	if err != nil {
		return returncodes.EACCES, errors.New("could not access " + fileName)
	}
	return returncodes.OK, nil
}
