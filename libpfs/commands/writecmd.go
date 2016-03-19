package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/encryption"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"io"
	"os"
	"path"
	"syscall"
)

//WriteCommand writes data to the given file
func WriteCommand(paranoidDirectory, filePath string, offset, length int64, data []byte) (returnCode returncodes.Code, returnError error, bytesWrote int) {
	Log.Info("write command called")
	Log.Verbose("write : given paranoidDirectory =", paranoidDirectory)

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

	code, err = canAccessFile(paranoidDirectory, inodeName, getAccessMode(syscall.O_WRONLY))
	if err != nil {
		return code, fmt.Errorf("unable to access %s: %s", filePath, err), 0
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
	contentsFile, err := os.OpenFile(path.Join(paranoidDirectory, "contents", inodeName), os.O_RDWR, 0777)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error opening contents file: %s", err), 0
	}
	defer contentsFile.Close()

	if offset == -1 {
		offset = 0
	}

	if length == -1 {
		err = truncate(contentsFile, offset)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error truncating file: %s", err), 0
		}
	} else if len(data) > int(length) {
		data = data[:length]
	} else if len(data) < int(length) {
		emptyBytes := make([]byte, int(length)-len(data))
		data = append(data, emptyBytes...)
	}

	wroteLen, err := writeAt(contentsFile, data, offset)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing to file: %s", err), wroteLen
	}

	return returncodes.OK, nil, wroteLen
}

func writeAt(file *os.File, data []byte, offset int64) (wroteLen int, err error) {
	if !encryption.Encrypted {
		return file.WriteAt(data, offset)
	}

	cipherSizeInt64 := int64(encryption.GetCipherSize())
	extraStartBytes := offset % cipherSizeInt64
	writeStart := offset - extraStartBytes
	startBytes := make([]byte, extraStartBytes)

	_, readerror, err := readAt(file, startBytes, writeStart)
	if err != nil {
		return 0, fmt.Errorf("error reading start block: %s", err)
	}
	if readerror != nil {
		return 0, fmt.Errorf("error reading start block: %s", readerror)
	}

	extraEndBytes := cipherSizeInt64 - ((offset + int64(len(data))) % cipherSizeInt64)
	endBytes := make([]byte, extraEndBytes)
	fileLength, err := getFileLength(file)
	if err != nil {
		return 0, err
	}

	if offset+int64(len(data)) < fileLength {
		_, readerror, err := readAt(file, endBytes, offset+int64(len(data)))
		if err != nil {
			return 0, fmt.Errorf("error reading end block: %s", err)
		}
		if readerror != nil && readerror != io.EOF {
			return 0, fmt.Errorf("error reading end block: %s", err)
		}
	}

	bytesToWrite := append(startBytes, data...)
	bytesToWrite = append(bytesToWrite, endBytes...)

	Log.Info("startBytes :", string(startBytes))
	Log.Info("data :", string(data))
	Log.Info("Bytes to write :", string(bytesToWrite))
	encBytes, err := encryption.Encrypt(bytesToWrite)
	if err != nil {
		return 0, err
	}

	encData := encBytes.Bytes()
	n, err := file.WriteAt(encData, writeStart+1)
	if err != nil {
		return n, err
	}

	Log.Info("offset : ", offset)
	Log.Info("data:", string(data))
	Log.Info("file length:", fileLength)
	if offset+int64(len(data)) > fileLength {
		endBlockSize := (offset + int64(len(data))) % cipherSizeInt64
		Log.Info("End block size : ", endBlockSize)
		_, err := file.WriteAt([]byte{byte(endBlockSize)}, 0)
		if err != nil {
			return n, err
		}
	}

	return n, nil
}
