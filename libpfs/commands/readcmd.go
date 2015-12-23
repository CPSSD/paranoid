package commands

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"io"
	"os"
	"path"
	"strconv"
	"syscall"
)

//ReadCommand reads data from a file
//Offset and length can be given as -1 to note that the defaults should be used.
func ReadCommand(directory, fileName string, offset, length int64) (returnCode int, returnError error, fileContents []byte) {
	Log.Info("read command called")
	Log.Verbose("read : given directory = " + directory)

	namepath := getParanoidPath(directory, fileName)

	err := getFileSystemLock(directory, sharedLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, nil
	}

	defer func() {
		err := unLockFileSystem(directory)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
			fileContents = nil
		}
	}()

	fileType, err := getFileType(directory, namepath)
	if err != nil {
		return returncodes.EUNEXPECTED, err, nil
	}

	if fileType == typeENOENT {
		return returncodes.ENOENT, errors.New(fileName + " does not exist"), nil
	}

	if fileType == typeDir {
		return returncodes.EISDIR, errors.New(fileName + " is a directory"), nil
	}

	if fileType == typeSymlink {
		return returncodes.EIO, errors.New(fileName + " is a symlink"), nil
	}

	inodeBytes, code, err := getFileInode(namepath)
	if code != returncodes.OK || err != nil {
		return code, err, nil
	}
	inodeFileName := string(inodeBytes)

	err = syscall.Access(path.Join(directory, "contents", inodeFileName), getAccessMode(syscall.O_RDONLY))
	if err != nil {
		return returncodes.EACCES, errors.New("could not access file " + fileName), nil
	}

	err = getFileLock(directory, inodeFileName, sharedLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, nil
	}

	defer func() {
		err := unLockFile(directory, inodeFileName)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
			fileContents = nil
		}
	}()

	file, err := os.OpenFile(path.Join(directory, "contents", inodeFileName), os.O_RDONLY, 0777)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error opening contents file", err), nil
	}

	var fileBuffer bytes.Buffer
	bytesRead := make([]byte, 1024)
	maxRead := 100000000

	if offset == -1 {
		offset = 0
	}

	if length != -1 {
		Log.Verbose("read : " + strconv.FormatInt(length, 10) + " bytes starting at " + strconv.FormatInt(offset, 10))
		maxRead = int(length)
	} else {
		Log.Verbose("read : from " + strconv.FormatInt(offset, 10) + " to end of file")
	}

	for {
		n, err := file.ReadAt(bytesRead, offset)
		if n > maxRead {
			bytesRead = bytesRead[0:maxRead]
			_, err := fileBuffer.Write(bytesRead)
			if err != nil {
				return returncodes.EUNEXPECTED, fmt.Errorf("error writing to file buffer:", err), nil
			}
			break
		}

		offset = offset + int64(n)
		maxRead = maxRead - n
		if err == io.EOF {
			bytesRead = bytesRead[:n]
			_, err := fileBuffer.Write(bytesRead)
			if err != nil {
				return returncodes.EUNEXPECTED, fmt.Errorf("error writing to file buffer:", err), nil
			}
			break
		}

		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error reading from "+fileName+":", err), nil
		}

		bytesRead = bytesRead[:n]
		_, err = fileBuffer.Write(bytesRead)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error writing to file buffer:", err), nil
		}
	}
	return returncodes.OK, nil, fileBuffer.Bytes()
}
