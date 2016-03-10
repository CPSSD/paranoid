package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"os"
	"path"
	"syscall"
)

//WriteCommand writes data to the given file
//offset and length can be given as -1 if the defaults are to be used
func WriteCommand(paranoidDirectory, filePath string, offset, length int64, data []byte) (returnCode int, returnError error, bytesWrote int) {
	Log.Verbose("write : given paranoidDirectory = " + paranoidDirectory)

	err := getFileSystemLock(paranoidDirectory, sharedLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, 0
	}

	defer func() {
		err := unLockFileSystem(paranoidDirectory)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
			bytesWrote = 0
		}
	}()

	namepath := getParanoidPath(paranoidDirectory, filePath)
	namepathType, err := getFileType(paranoidDirectory, namepath)
	if err != nil {
		return returncodes.EUNEXPECTED, err, 0
	}

	if namepathType == typeENOENT {
		return returncodes.ENOENT, errors.New(filePath + " does not exist"), 0
	}

	if namepathType == typeDir {
		return returncodes.EISDIR, errors.New(filePath + " is a paranoidDirectory"), 0
	}

	if namepathType == typeSymlink {
		return returncodes.EIO, errors.New(filePath + " is a symlink"), 0
	}

	fileInodeBytes, code, err := getFileInode(namepath)
	if code != returncodes.OK {
		return code, err, 0
	}
	inodeName := string(fileInodeBytes)

	err = syscall.Access(path.Join(paranoidDirectory, "contents", inodeName), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		return returncodes.EACCES, errors.New("could not access " + filePath), 0
	}

	err = getFileLock(paranoidDirectory, inodeName, exclusiveLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, 0
	}

	defer func() {
		err := unLockFile(paranoidDirectory, inodeName)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
			bytesWrote = 0
		}
	}()

	Log.Verbose("write : wrting to " + inodeName)
	contentsFile, err := os.OpenFile(path.Join(paranoidDirectory, "contents", inodeName), os.O_WRONLY, 0777)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error opening contents file:", err), 0
	}
	defer contentsFile.Close()

	if offset == -1 {
		offset = 0
	}

	if length == -1 {
		err = contentsFile.Truncate(offset)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error truncating file:", err), 0
		}
	} else {
		if len(data) > int(length) {
			data = data[:length]
		} else if len(data) < int(length) {
			emptyBytes := make([]byte, int(length)-len(data))
			data = append(data, emptyBytes...)
		}
	}

	wroteLen, err := contentsFile.WriteAt(data, offset)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing to file:", err), wroteLen
	}

	return returncodes.OK, nil, wroteLen
}
