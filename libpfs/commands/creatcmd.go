package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/encryption"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"io/ioutil"
	"os"
	"path"
)

//CreatCommand creates a new file with the name filePath in the pfs paranoidDirectory
func CreatCommand(paranoidDirectory, filePath string, perms os.FileMode, shouldGlob bool) (returnCode returncodes.Code, returnError error) {
	Log.Info("creat command called")
	Log.Verbose("creat : paranoidDirectory = " + paranoidDirectory)

	err := GetFileSystemLock(paranoidDirectory, ExclusiveLock)
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

	namepath := GetParanoidPath(paranoidDirectory, filePath)

	fileType, err := getFileType(paranoidDirectory, namepath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	if fileType != typeENOENT {
		return returncodes.EEXIST, errors.New(filePath + " already exists")
	}
	Log.Verbose("creat : creating file " + filePath)

	uuidbytes, err := generateNewInode()
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	uuidstring := string(uuidbytes)
	Log.Verbose("creat : uuid = " + uuidstring)

	err = ioutil.WriteFile(namepath, uuidbytes, 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, errors.New("error writing name file")
	}

	nodeData := &Inode{
		Mode:    perms,
		Inode:   uuidstring,
		Count:   1,
		Ignored: shouldGlob}

	jsonData, err := json.Marshal(nodeData)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error marshalling json: %s", err)
	}

	err = ioutil.WriteFile(path.Join(paranoidDirectory, "inodes", uuidstring), jsonData, 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing inodes file: %s", err)
	}

	contentsFile, err := os.Create(path.Join(paranoidDirectory, "contents", uuidstring))
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error creating contents file: %s", err)
	}
	defer contentsFile.Close()

	if encryption.Encrypted {
		n, err := contentsFile.WriteAt([]byte{1}, 0)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error creating contents file: %s", err)
		}
		if n != 1 {
			return returncodes.EUNEXPECTED, errors.New("error writing first byte to contents file")
		}
	}

	return returncodes.OK, nil
}
