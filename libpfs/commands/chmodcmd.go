package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"os"
	"path"
	"syscall"
)

//ChmodCommand is used to change the permissions of a file.
func ChmodCommand(paranoidDirectory, filePath string, perms os.FileMode) (returnCode int, returnError error) {
	Log.Info("chmod command called")
	Log.Verbose("chmod : given paranoidDirectory = " + paranoidDirectory)

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

	if fileType == typeENOENT {
		return returncodes.ENOENT, errors.New(filePath + " does not exist")
	}

	if fileType == typeSymlink {
		return returncodes.EIO, errors.New(filePath + " is of type symlink")
	}

	inodeNameBytes, code, err := getFileInode(namepath)
	if code != returncodes.OK {
		return code, err
	}
	inodeName := string(inodeNameBytes)

	err = syscall.Access(path.Join(paranoidDirectory, "contents", inodeName), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		return returncodes.EACCES, errors.New("unable to access " + filePath)
	}

	Log.Verbosef("chmod : changing permissions of "+inodeName+" to", perms)

	contentsFile, err := os.OpenFile(path.Join(paranoidDirectory, "contents", inodeName), os.O_WRONLY, 0777)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("unexpected error attempting to open file:", err)
	}
	defer contentsFile.Close()

	err = contentsFile.Chmod(perms)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("unexpected error attempting to change file permissions:", err)
	}

	return returncodes.OK, nil
}
