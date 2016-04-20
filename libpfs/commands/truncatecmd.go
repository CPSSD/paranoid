package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/encryption"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"os"
	"path"
)

//TruncateCommand reduces the file given to the new length
func TruncateCommand(paranoidDirectory, filePath string, length int64) (returnCode returncodes.Code, returnError error) {
	Log.Info("truncate command called")
	Log.Verbose("truncate : given paranoidDirectory = " + paranoidDirectory)

	err := GetFileSystemLock(paranoidDirectory, SharedLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	defer func() {
		err := UnLockFileSystem(paranoidDirectory)
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

	fileInodeBytes, code, err := GetFileInode(namepath)
	if code != returncodes.OK {
		return code, err
	}
	inodeName := string(fileInodeBytes)

	err = getFileLock(paranoidDirectory, inodeName, ExclusiveLock)
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

	contentsFile, err := os.OpenFile(path.Join(paranoidDirectory, "contents", inodeName), os.O_WRONLY, 0777)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error opening contents file: %s", err)
	}
	defer contentsFile.Close()

	err = truncate(contentsFile, length)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error truncating file: %s", err)
	}

	return returncodes.OK, nil
}

func truncate(file *os.File, length int64) error {
	if !encryption.Encrypted {
		return file.Truncate(length)
	}

	cipherSizeInt64 := int64(encryption.GetCipherSize())
	newLastBlockLength := length % cipherSizeInt64

	if newLastBlockLength == 0 {
		err := file.Truncate(length + 1)
		if err != nil {
			return err
		}
		_, err = file.WriteAt([]byte{byte(cipherSizeInt64)}, 0)
		if err != nil {
			return err
		}
		return nil
	}

	err := file.Truncate(length + 1 + cipherSizeInt64 - newLastBlockLength)
	if err != nil {
		return err
	}

	_, err = file.WriteAt([]byte{byte(newLastBlockLength)}, 0)
	if err != nil {
		return err
	}
	return nil
}
