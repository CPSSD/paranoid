package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/pfsd/pnetclient"
	"os"
	"path"
	"syscall"
)

//ChmodCommand is used to change the permissions of a file.
func ChmodCommand(directory, fileName string, perms os.FileMode, sendOverNetwork bool) (returnCode int, returnError error) {
	Log.Info("chmod command given")
	Log.Verbose("chmod : given directory = " + directory)

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

	namepath := getParanoidPath(directory, fileName)

	fileType, err := getFileType(directory, namepath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	if fileType == typeENOENT {
		return returncodes.ENOENT, errors.New(fileName + " does not exist")
	}

	if fileType == typeSymlink {
		return returncodes.EIO, errors.New(fileName + " is of type symlink")
	}

	inodeNameBytes, code, err := getFileInode(namepath)
	if code != returncodes.OK {
		return code, err
	}
	inodeName := string(inodeNameBytes)

	err = syscall.Access(path.Join(directory, "contents", inodeName), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		return returncodes.EACCES, errors.New("unable to access " + fileName)
	}

	Log.Verbosef("chmod : changing permissions of "+inodeName+" to", perms)

	contentsFile, err := os.OpenFile(path.Join(directory, "contents", inodeName), os.O_WRONLY, 0777)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("unexpected error attempting to open file:", err)
	}

	err = contentsFile.Chmod(perms)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("unexpected error attempting to change file permissions:", err)
	}

	if sendOverNetwork {
		pnetclient.Chmod(fileName, uint32(perms))
	}
	return returncodes.OK, nil
}
