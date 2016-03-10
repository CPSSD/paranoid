package commands

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/pfsd/globals"
	"io"
	"os"
	"path"
	"syscall"
)

// ReadCommand reads data from a file
// Offset and length can be given as -1 to note that the defaults should be used.
func ReadCommand(paranoidDirectory, filePath string, offset, length int64) (returnCode int, returnError error, fileContents []byte) {
	libpfs.CipherBlock, err = libpfs.GenerateAESCipherBlock(globals.EncryptionKey.GetBytes())
	if err != nil {
		Log.Fatalf("failed to initialize cipher block: %s", err)
	}
	if !libpfs.CheckIsCipher() {
		Log.Fatal("cipherBlock did not initialize correctly")
	}

	Log.Info("read command called")
	Log.Verbose("read : given paranoidDirectory = " + paranoidDirectory)

	namepath := getParanoidPath(paranoidDirectory, filePath)

	err := getFileSystemLock(paranoidDirectory, sharedLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, nil
	}

	defer func() {
		err := unLockFileSystem(paranoidDirectory)
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

	inodeBytes, code, err := getFileInode(namepath)
	if code != returncodes.OK || err != nil {
		return code, err, nil
	}
	inodeFileName := string(inodeBytes)

	code, err = canAccessFile(paranoidDirectory, inodeFileName, getAccessMode(syscall.O_RDONLY))
	if err != nil {
		return code, fmt.Errorf("unable to access %s: %s", filePath, err), nil
	}

	err = getFileLock(paranoidDirectory, inodeFileName, sharedLock)
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

	// Read the file and return if successful
	data, err := read(file, offset, length)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error reading: %s", err), nil
	}
	return returncodes.OK, nil, data.Bytes()
}

// read is a function only for reading the encrypted file
func read(file *os.File, offset int64, length int64) (decBuf *bytes.Buffer, err error) {
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
		return nil, err
	}
	fileSize := stats.Size() - 1 // Take the correction for the offset

	// Check the size of the last block
	lastBlockSize, err := libpfs.LastBlockSize(file)
	if err != nil {
		return nil, err
	}

	// Check the length size
	if length == -1 {
		length = fileSize - (blockSize - int64(lastBlockSize))
	}

	// Define the block length
	blockLength := length - length%blockSize + 1 + blockSize

	// Read the file
	buf := make([]byte, blockLength-blockOffset)
	_, err = file.ReadAt(buf, blockOffset)
	if err != nil && err != io.EOF {
		return nil, err
	}

	// Decrypt the file
	dec := libpfs.Decrypt(buf, lastBlockSize)
	dec.Next(int(offset - blockOffset))
	dec.Truncate(int(length))

	return decBuf, nil
}
