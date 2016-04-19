package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"io/ioutil"
	"os"
)

// RenameCommand is called when renaming a file
func RenameCommand(paranoidDirectory, oldFilePath, newFilePath string) (returnCode returncodes.Code, returnError error) {
	Log.Info("rename command called")
	oldFileParanoidPath := GetParanoidPath(paranoidDirectory, oldFilePath)
	newFileParanoidPath := GetParanoidPath(paranoidDirectory, newFilePath)

	oldFileType, err := getFileType(paranoidDirectory, oldFileParanoidPath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	newFileType, err := getFileType(paranoidDirectory, newFileParanoidPath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	err = GetFileSystemLock(paranoidDirectory, ExclusiveLock)
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

	if oldFileType == typeENOENT {
		return returncodes.ENOENT, errors.New(oldFilePath + " does not exist")
	}

	if newFileType != typeENOENT {
		//Renaming is allowed when a file already exists, unless the existing file is a non empty paranoidDirectory
		if newFileType == typeFile {
			_, err := UnlinkCommand(paranoidDirectory, newFilePath)
			if err != nil {
				return returncodes.EEXIST, errors.New(newFilePath + " already exists")
			}
		} else if newFileType == typeDir {
			dirpath := GetParanoidPath(paranoidDirectory, newFilePath)
			files, err := ioutil.ReadDir(dirpath)
			if err != nil || len(files) > 0 {
				return returncodes.ENOTEMPTY, errors.New(newFilePath + " is not empty")
			}
			_, err = RmdirCommand(paranoidDirectory, newFilePath)
			if err != nil {
				return returncodes.EEXIST, errors.New(newFilePath + " already exists")
			}
		}
	}

	err = os.Rename(oldFileParanoidPath, newFileParanoidPath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error renaming file: %s", err)
	}

	return returncodes.OK, nil
}
