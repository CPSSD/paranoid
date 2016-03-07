package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"os"
	"path"
	"syscall"
	"time"
)

//UtimesCommand updates the acess time and modified time of a file
func UtimesCommand(paranoidDirectory, filePath string, atime, mtime *time.Time) (returnCode int, returnError error) {
	Log.Info("utimes command called")
	Log.Verbose("utimes : given paranoidDirectory = " + paranoidDirectory)

	err := getFileSystemLock(paranoidDirectory, sharedLock)
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

	fileInodeBytes, code, err := getFileInode(namepath)
	if code != returncodes.OK {
		return code, err
	}
	inodeName := string(fileInodeBytes)

	err = syscall.Access(path.Join(paranoidDirectory, "contents", inodeName), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		return returncodes.EACCES, errors.New("could not access " + filePath)
	}

	err = getFileLock(paranoidDirectory, inodeName, exclusiveLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	defer func() {
		err := unLockFile(paranoidDirectory, inodeName)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
		}
	}()

	file, err := os.Open(path.Join(paranoidDirectory, "contents", inodeName))
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error opening contents file:", err)
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error stating file:", err)
	}

	stat := fi.Sys().(*syscall.Stat_t)
	oldatime := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
	oldmtime := fi.ModTime()
	if atime == nil && mtime == nil {
		return returncodes.EUNEXPECTED, errors.New("no times to update!")
	}

	if atime == nil {
		err = os.Chtimes(path.Join(paranoidDirectory, "contents", inodeName), oldatime, *mtime)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error changing times:", err)
		}
	} else if mtime == nil {
		err = os.Chtimes(path.Join(paranoidDirectory, "contents", inodeName), *atime, oldmtime)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error changing times:", err)
		}
	} else {
		err = os.Chtimes(path.Join(paranoidDirectory, "contents", inodeName), *atime, *mtime)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error changing times:", err)
		}
	}

	return returncodes.OK, nil
}
