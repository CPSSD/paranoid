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

//TruncateCommand reduces the file given to the new length
func TruncateCommand(directory, fileName string, length int64, sendOverNetwork bool) (returnCode int, returnError error) {
	Log.Info("truncate command given")
	Log.Verbose("truncate : given directory = " + directory)

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

	namepath := getParanoidPath(directory, fileName)
	namepathType, err := getFileType(directory, namepath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	if namepathType == typeENOENT {
		return returncodes.ENOENT, errors.New(fileName + " does not exist")
	}

	if namepathType == typeDir {
		return returncodes.EISDIR, errors.New(fileName + " is a directory")
	}

	if namepathType == typeSymlink {
		return returncodes.EIO, errors.New(fileName + " is a symlink")
	}

	fileInodeBytes, code, err := getFileInode(namepath)
	if code != returncodes.OK {
		return code, err
	}
	inodeName := string(fileInodeBytes)

	err = syscall.Access(path.Join(directory, "contents", inodeName), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		return returncodes.EACCES, errors.New("could not access " + fileName)
	}

	err = getFileLock(directory, inodeName, exclusiveLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	defer func() {
		err := unLockFile(directory, inodeName)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
		}
	}()

	Log.Verbose("truncate : truncating " + fileName)

	contentsFile, err := os.OpenFile(path.Join(directory, "contents", inodeName), os.O_WRONLY, 0777)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error opening contents file:", err)
	}
	defer contentsFile.Close()

	err = contentsFile.Truncate(length)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error truncating file:", err)
	}

	if sendOverNetwork {
		pnetclient.Truncate(fileName, uint64(length))
	}
	return returncodes.OK, nil
}
