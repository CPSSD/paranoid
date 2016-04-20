package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"io/ioutil"
	"path"
)

// ReadlinkCommand reads the value of the symbolic link
func ReadlinkCommand(paranoidDirectory, filePath string) (returnCode returncodes.Code, returnError error, linkContents string) {
	Log.Info("readlink command called")

	err := GetFileSystemLock(paranoidDirectory, SharedLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, ""
	}

	defer func() {
		err := UnLockFileSystem(paranoidDirectory)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
			linkContents = ""
		}
	}()

	link := getParanoidPath(paranoidDirectory, filePath)
	fileType, err := getFileType(paranoidDirectory, link)
	if err != nil {
		return returncodes.EUNEXPECTED, err, ""
	}

	if fileType == typeENOENT {
		return returncodes.ENOENT, errors.New(filePath + " does not exist"), ""
	}

	if fileType == typeDir {
		return returncodes.EISDIR, errors.New(filePath + " is a paranoidDirectory"), ""
	}

	if fileType == typeFile {
		return returncodes.EIO, errors.New(filePath + " is a file"), ""
	}

	Log.Verbose("readlink: given paranoidDirectory", paranoidDirectory)

	linkInode, code, err := GetFileInode(link)
	if code != returncodes.OK || err != nil {
		return code, err, ""
	}

	err = getFileLock(paranoidDirectory, string(linkInode), SharedLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, ""
	}

	defer func() {
		err := unLockFile(paranoidDirectory, string(linkInode))
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
			linkContents = ""
		}
	}()

	inodePath := path.Join(paranoidDirectory, "inodes", string(linkInode))

	inodeContents, err := ioutil.ReadFile(inodePath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error reading link: %s", err), ""
	}

	inodeData := &inode{}
	Log.Verbose("readlink unmarshaling ", string(inodeContents))
	err = json.Unmarshal(inodeContents, &inodeData)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error unmarshalling json: %s", err), ""
	}

	return returncodes.OK, nil, inodeData.Link
}
