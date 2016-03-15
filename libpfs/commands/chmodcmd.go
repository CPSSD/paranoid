package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"io/ioutil"
	"os"
	"path"
	"syscall"
)

const PERM_MASK = 0777

//ChmodCommand is used to change the permissions of a file.
func ChmodCommand(paranoidDirectory, filePath string, perms os.FileMode) (returnCode int, returnError error) {
	Log.Info("chmod command called")
	Log.Verbose("chmod : given paranoidDirectory = " + paranoidDirectory)

	err := getFileSystemLock(paranoidDirectory, exclusiveLock)
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

	fileType, err := getFileType(paranoidDirectory, namepath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	if fileType == typeENOENT {
		return returncodes.ENOENT, errors.New(filePath + " does not exist")
	}

	if fileType == typeSymlink {
		return returncodes.EIO, errors.New(filePath + " is of type symlink")
	}

	inodeNameBytes, code, err := getFileInode(namepath)
	if code != returncodes.OK {
		return code, err
	}
	inodeName := string(inodeNameBytes)

	code, err = canAccessFile(paranoidDirectory, inodeName, getAccessMode(syscall.O_WRONLY))
	if err != nil {
		return code, fmt.Errorf("unable to access %s: %s", filePath, err)
	}

	Log.Verbosef("chmod : changing permissions of "+inodeName+" to", perms)

	inodePath := path.Join(paranoidDirectory, "inodes", inodeName)
	inodeContents, err := ioutil.ReadFile(inodePath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error reading inode: %s", err)
	}

	nodeData := &inode{}
	err = json.Unmarshal(inodeContents, &nodeData)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error unmarshaling inode data: %s", err)
	}

	nodeData.Mode = (nodeData.Mode &^ PERM_MASK) | perms

	jsonData, err := json.Marshal(nodeData)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error marshalling json:", err)
	}

	err = ioutil.WriteFile(inodePath, jsonData, 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing inodes file:", err)
	}

	return returncodes.OK, nil
}
