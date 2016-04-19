package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"io/ioutil"
	"os"
	"path"
)

// UnlinkCommand removes a filename link from an inode.
func UnlinkCommand(paranoidDirectory, filePath string) (returnCode returncodes.Code, returnError error) {
	Log.Info("unlink command called")
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

	fileParanoidPath := GetParanoidPath(paranoidDirectory, filePath)
	fileType, err := getFileType(paranoidDirectory, fileParanoidPath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	Log.Verbose("unlink : paranoidDirectory given = " + paranoidDirectory)

	// checking if file exists
	if fileType == typeENOENT {
		return returncodes.ENOENT, errors.New(filePath + " does not exist")
	}

	if fileType == typeDir {
		return returncodes.EISDIR, errors.New(filePath + " is a paranoidDirectory")
	}

	// getting file inode
	inodeBytes, code, err := GetFileInode(fileParanoidPath)
	if code != returncodes.OK {
		return code, err
	}

	// removing filename
	Log.Verbose("unlink : deleting file " + fileParanoidPath)
	err = os.Remove(fileParanoidPath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error removing file in names: %s", err)
	}

	// getting inode contents
	inodePath := path.Join(paranoidDirectory, "inodes", string(inodeBytes))
	Log.Verbose("unlink : reading file " + inodePath)
	inodeContents, err := ioutil.ReadFile(inodePath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error reading inodes contents: %s", err)
	}

	inodeData := &Inode{}
	Log.Verbose("unlink unmarshaling ", string(inodeContents))
	err = json.Unmarshal(inodeContents, &inodeData)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error unmarshaling json: %s", err)
	}

	if inodeData.Count == 1 {
		// remove inode and contents
		contentsPath := path.Join(paranoidDirectory, "contents", string(inodeBytes))
		Log.Verbose("unlink : removing file " + contentsPath)
		err = os.Remove(contentsPath)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error removing contents: %s", err)
		}

		Log.Verbose("unlink : removing file " + inodePath)
		err = os.Remove(inodePath)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error removing inode: %s", err)
		}
	} else {
		// subtracting one from inode count and saving
		inodeData.Count--
		Log.Verbose("unlink : truncating file " + inodePath)
		err = os.Truncate(inodePath, 0)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error truncating inode path: %s", err)
		}

		dataToWrite, err := json.Marshal(inodeData)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error marshalling json: %s", err)
		}

		Log.Verbose("unlink : writing to file " + inodePath)
		err = ioutil.WriteFile(inodePath, dataToWrite, 0777)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error writing to inode file: %s", err)
		}
	}

	return returncodes.OK, nil
}
