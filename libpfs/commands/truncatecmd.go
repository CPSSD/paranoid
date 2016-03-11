package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"os"
	"path"
	"syscall"
)

//TruncateCommand reduces the file given to the new length
func TruncateCommand(paranoidDirectory, filePath string, length int64) (returnCode returncodes.Code, returnError error) {
	Log.Info("truncate command called")
	Log.Verbose("truncate : given paranoidDirectory = " + paranoidDirectory)

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
	namepathType, err := getFileType(paranoidDirectory, namepath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	if namepathType == typeENOENT {
		return returncodes.ENOENT, errors.New(filePath + " does not exist")
	}

	if namepathType == typeDir {
		return returncodes.EISDIR, errors.New(filePath + " is a paranoidDirectory")
	}

	if namepathType == typeSymlink {
		return returncodes.EIO, errors.New(filePath + " is a symlink")
	}

	fileInodeBytes, code, err := getFileInode(namepath)
	if code != returncodes.OK {
		return code, err
	}
	inodeName := string(fileInodeBytes)

	code, err = canAccessFile(paranoidDirectory, inodeName, getAccessMode(syscall.O_WRONLY))
	if err != nil {
		return code, fmt.Errorf("unable to access %s: %s", filePath, err)
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

	Log.Verbose("truncate : truncating " + filePath)

	contentsFile, err := os.OpenFile(path.Join(paranoidDirectory, "contents", inodeName), os.O_RDWR, 0777)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error opening contents file: %s", err)
	}
	defer contentsFile.Close()

	err = truncate(contentsFile, length)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("unable to trucate: %s", err)
	}

	return returncodes.OK, nil
}

func truncate(file *os.File, len int64) (err error) {
	blockSize := int64(libpfs.CipherBlock.BlockSize())

	// offset of the last block
	blockOffset := len - len%blockSize + 1

	finalBlock, err := libpfs.GetLastBlock(file)
	if err != nil {
		return fmt.Errorf("can't read last block: %s", err)
	}

	// truncate to the size blockOffset
	err = file.Truncate(blockOffset)
	if err != nil {
		return fmt.Errorf("unable to truncate at block offset: %s", err)
	}

	// Get the size of the last block
	l, err := libpfs.LastBlockSize(file)
	if err != nil {
		return fmt.Errorf("unable to get the last block size: %s", err)
	}

	// Decode the block and cut it to size
	dec := libpfs.Decrypt(finalBlock, l)
	finalBlock = dec.Bytes()[:len%blockSize]

	// Encrypt the last file again and save at the end
	enc, l := libpfs.Encrypt(finalBlock)

	_, err = file.WriteAt([]byte{byte(l)}, 0)
	if err != nil {
		return fmt.Errorf("nable to write the last block size: %s", err)
	}
	_, err = file.WriteAt(enc.Bytes(), blockOffset)
	if err != nil {
		return fmt.Errorf("unable to write re-encrypted data: %s", err)
	}

	return nil
}
