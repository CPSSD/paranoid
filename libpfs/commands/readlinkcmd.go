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
func ReadlinkCommand(directory, fileName string) (returnCode int, returnError error, linkContents string) {
	Log.Info("readlink called")

	err := getFileSystemLock(directory, sharedLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, ""
	}

	defer func() {
		err := unLockFileSystem(directory)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
			linkContents = ""
		}
	}()

	link := getParanoidPath(directory, fileName)
	fileType, err := getFileType(directory, link)
	if err != nil {
		return returncodes.EUNEXPECTED, err, ""
	}

	if fileType == typeENOENT {
		return returncodes.ENOENT, errors.New(fileName + " does not exist"), ""
	}

	if fileType == typeDir {
		return returncodes.EISDIR, errors.New(fileName + " is a directory"), ""
	}

	if fileType == typeFile {
		return returncodes.EIO, errors.New(fileName + " is a file"), ""
	}

	Log.Verbose("readlink: given directory", directory)

	linkInode, code, err := getFileInode(link)
	if code != returncodes.OK || err != nil {
		return code, err, ""
	}

	err = getFileLock(directory, string(linkInode), sharedLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, ""
	}

	defer func() {
		err := unLockFile(directory, string(linkInode))
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
			linkContents = ""
		}
	}()

	inodePath := path.Join(directory, "inodes", string(linkInode))

	inodeContents, err := ioutil.ReadFile(inodePath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error reading link:", err), ""
	}

	inodeData := &inode{}
	Log.Verbose("readlink unmarshaling ", string(inodeContents))
	err = json.Unmarshal(inodeContents, &inodeData)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error unmarshalling json:", err), ""
	}

	return returncodes.OK, nil, inodeData.Link
}
