package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	// "io"
	"os"
	"path"
	"syscall"
)

//WriteCommand writes data to the given file
//offset and length can be given as -1 if the defaults are to be used
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

	switch namepathType {
	case typeENOENT:
		return returncodes.ENOENT, fmt.Errorf("%s does not exist", filePath), 0
	case typeDir:
		return returncodes.EISDIR, fmt.Errorf("%s is a paranoidDirectory", filePath), 0
	case typeSymlink:
		return returncodes.EIO, fmt.Errorf("%s is a symlink", filePath), 0
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

	Log.Verbose("write : writing to", inodeName)
	contentsFile, err := os.OpenFile(path.Join(paranoidDirectory, "contents", inodeName), os.O_RDWR, 0777)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error opening contents file: %s", err), 0
	}
	defer contentsFile.Close()

	// write the data to file
	wroteLen, err := write(contentsFile, data, offset, length)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing data to file: %s", err), 0
	}
	return returncodes.OK, nil, wroteLen
}

func write(file *os.File, data []byte, offset int64, length int64) (n int, err error) {
	blockSize := int64(libpfs.CipherBlock.BlockSize())

	// Check the offset
	if offset <= 0 {
		offset = 0
	}

	// Define block offset to read from
	blockOffset := offset - offset%blockSize + 1

	// Get the total size of the file
	stats, err := file.Stat()
	if err != nil {
		return
	}
	fileSize := stats.Size() - 1 // Take the correction for the offset

	// Check the size of the last block
	lastBlockSize, err := libpfs.LastBlockSize(file)
	if err != nil {
		return 0, fmt.Errorf("unable to read last block size: %s", err)
	}

	// Check the length size. Truncate if length == -1
	if length == -1 {
		err = truncate(file, offset)
	} else {
		// Append empty bytes if length > len(data)
		if length > int64(len(data)) {
			data = append(data, make([]byte, len(data)-int(length))...)
			length = int64(len(data))
		} else {
			data = data[:length]
		}
	}

	// Define the block length (end of full block after wanted length)
	// blockLength := length - length%blockSize + 1 + blockSize

	// Read the last block and add it on if the last block is not full
	if lastBlockSize != 0 {
		lastBlock, err := libpfs.GetLastBlock(file, fileSize)
		if err != nil {
			return 0, fmt.Errorf("error getting last block: %s", err)
		}

		dec := libpfs.Decrypt(lastBlock, lastBlockSize)
		data = append(dec.Bytes(), data...)
	}

	// Encrypt the data
	enc, l := libpfs.Encrypt(data)

	// Write the encrypted block size at the beginning of the file
	_, err = file.WriteAt([]byte{byte(l)}, 0)
	if err != nil {
		return 0, fmt.Errorf("cannot write last block size: %s", err)
	}

	// Write the encrypted data to the file
	_, err = file.WriteAt(enc.Bytes(), blockOffset)
	if err != nil {
		return 0, fmt.Errorf("unable to write encrypted data: %s", err)
	}

	return int(length), nil
}
