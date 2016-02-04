package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/pfsd/pnetclient"
	"os"
	"path"
	"syscall"
	"time"
)

//UtimesCommand updates the acess time and modified time of a file
func UtimesCommand(directory, fileName string, atime, mtime *time.Time, sendOverNetwork bool) (returnCode int, returnError error) {
	Log.Info("utimes command called")
	Log.Verbose("utimes : given directory = " + directory)

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

	fileType, err := getFileType(directory, namepath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	if fileType == typeENOENT {
		return returncodes.ENOENT, errors.New(fileName + " does not exist")
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

	file, err := os.Open(path.Join(directory, "contents", inodeName))
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error opening contents file:", err)
	}

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
		err = os.Chtimes(path.Join(directory, "contents", inodeName), oldatime, *mtime)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error changing times:", err)
		}
		if sendOverNetwork {
			pnetclient.Utimes(fileName, 0, 0, int64(mtime.Second()), int64(mtime.Nanosecond()))
		}
	} else if mtime == nil {
		err = os.Chtimes(path.Join(directory, "contents", inodeName), *atime, oldmtime)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error changing times:", err)
		}
		if sendOverNetwork {
			pnetclient.Utimes(fileName, int64(atime.Second()), int64(atime.Nanosecond()), 0, 0)
		}
	} else {
		err = os.Chtimes(path.Join(directory, "contents", inodeName), *atime, *mtime)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error changing times:", err)
		}
		if sendOverNetwork {
			pnetclient.Utimes(fileName, int64(atime.Second()), int64(atime.Nanosecond()), int64(mtime.Second()), int64(mtime.Nanosecond()))
		}
	}

	return returncodes.OK, nil
}
