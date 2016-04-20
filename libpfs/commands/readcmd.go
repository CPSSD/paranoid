package commands

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/encryption"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"io"
	"os"
	"path"
	"strconv"
)

//ReadCommand reads data from a file
func ReadCommand(paranoidDirectory, filePath string, offset, length int64) (returnCode returncodes.Code, returnError error, fileContents []byte) {
	Log.Verbose("read : given paranoidDirectory = " + paranoidDirectory)

	namepath := getParanoidPath(paranoidDirectory, filePath)

	err := GetFileSystemLock(paranoidDirectory, SharedLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, nil
	}

	defer func() {
		err := UnLockFileSystem(paranoidDirectory)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
			fileContents = nil
		}
	}()

	fileType, err := getFileType(paranoidDirectory, namepath)
	if err != nil {
		return returncodes.EUNEXPECTED, err, nil
	}

	if fileType == typeENOENT {
		return returncodes.ENOENT, errors.New(filePath + " does not exist"), nil
	}

	if fileType == typeDir {
		return returncodes.EISDIR, errors.New(filePath + " is a paranoidDirectory"), nil
	}

	if fileType == typeSymlink {
		return returncodes.EIO, errors.New(filePath + " is a symlink"), nil
	}

	inodeBytes, code, err := GetFileInode(namepath)
	if code != returncodes.OK || err != nil {
		return code, err, nil
	}
	inodeFileName := string(inodeBytes)

	err = getFileLock(paranoidDirectory, inodeFileName, SharedLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, nil
	}

	defer func() {
		err := unLockFile(paranoidDirectory, inodeFileName)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
			fileContents = nil
		}
	}()

	file, err := os.OpenFile(path.Join(paranoidDirectory, "contents", inodeFileName), os.O_RDONLY, 0777)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error opening contents file: %s", err), nil
	}
	defer file.Close()

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
		n, readerror, err := readAt(file, bytesRead, offset)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error reading file %s", err), nil
		}

		if n > maxRead {
			bytesRead = bytesRead[0:maxRead]
			_, err := fileBuffer.Write(bytesRead)
			if err != nil {
				return returncodes.EUNEXPECTED, fmt.Errorf("error writing to file buffer: %s", err), nil
			}
			break
		}

		offset = offset + int64(n)
		maxRead = maxRead - n
		if readerror == io.EOF {
			bytesRead = bytesRead[:n]
			_, err := fileBuffer.Write(bytesRead)
			if err != nil {
				return returncodes.EUNEXPECTED, fmt.Errorf("error writing to file buffer: %s", err), nil
			}
			break
		}

		if readerror != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error reading from %s: %s", filePath, err), nil
		}

		bytesRead = bytesRead[:n]
		_, err = fileBuffer.Write(bytesRead)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error writing to file buffer: %s", err), nil
		}
	}
	return returncodes.OK, nil, fileBuffer.Bytes()
}

func readAt(file *os.File, bytesRead []byte, offset int64) (n int, readerror error, err error) {
	if !encryption.Encrypted {
		n, readerror := file.ReadAt(bytesRead, offset)
		return n, readerror, nil
	}

	if len(bytesRead) == 0 {
		return 0, nil, nil
	}

	cipherSizeInt64 := int64(encryption.GetCipherSize())
	extraStartBytes := offset % cipherSizeInt64
	extraEndBytes := cipherSizeInt64 - ((offset + int64(len(bytesRead))) % cipherSizeInt64)
	readStart := 1 + offset - extraStartBytes
	newBytesRead := make([]byte, int64(len(bytesRead))+extraStartBytes+extraEndBytes)

	fileLength, err := getFileLength(file)
	if err != nil {
		return 0, nil, err
	}

	n, readerror = file.ReadAt(newBytesRead, readStart)
	n = n - int(extraStartBytes)
	if n > len(bytesRead) {
		n = len(bytesRead)
	}
	if offset+int64(n) > fileLength {
		n = int(fileLength - offset)
	}

	err = encryption.Decrypt(newBytesRead)
	if err != nil {
		return 0, nil, err
	}
	newBytesRead = newBytesRead[extraStartBytes:]
	copy(bytesRead, newBytesRead)
	return n, readerror, nil
}
