package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"os"
	"path"
	"syscall"
)

//WriteCommand writes data to the given file
//offset and length can be given as -1 if the defaults are to be used
func WriteCommand(directory, fileName string, offset, length int64, data []byte, sendOverNetwork bool) (returnCode int, returnError error, bytesWrote int) {
	Log.Info("write command given")
	Log.Verbose("write : given directory = " + directory)

	err := getFileSystemLock(directory, sharedLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, 0
	}

	defer func() {
		err := unLockFileSystem(directory)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
			bytesWrote = 0
		}
	}()

	namepath := getParanoidPath(directory, fileName)
	namepathType, err := getFileType(directory, namepath)
	if err != nil {
		return returncodes.EUNEXPECTED, err, 0
	}

	if namepathType == typeENOENT {
		return returncodes.ENOENT, errors.New(fileName + " does not exist"), 0
	}

	if namepathType == typeDir {
		return returncodes.EISDIR, errors.New(fileName + " is a directory"), 0
	}

	if namepathType == typeSymlink {
		return returncodes.EIO, errors.New(fileName + " is a symlink"), 0
	}

	fileInodeBytes, code, err := getFileInode(namepath)
	if code != returncodes.OK {
		return code, err, 0
	}
	inodeName := string(fileInodeBytes)

	err = syscall.Access(path.Join(directory, "contents", inodeName), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		return returncodes.EACCES, errors.New("could not access " + fileName), 0
	}

	err = getFileLock(directory, inodeName, exclusiveLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, 0
	}

	defer func() {
		err := unLockFile(directory, inodeName)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
			bytesWrote = 0
		}
	}()

	Log.Verbose("write : wrting to " + inodeName)
	contentsFile, err := os.OpenFile(path.Join(directory, "contents", inodeName), os.O_WRONLY, 0777)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error opening contents file:", err), 0
	}

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

	if sendOverNetwork {
		//do this when mega binary
		//sendToServer(directory, "write", args[1:], fileData)
	}
	return returncodes.OK, nil, wroteLen
}
