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

// UnlinkCommand removes a filename link from an inode.
func UnlinkCommand(paranoidDirectory, filePath string) (returnCode int, returnError error) {
	Log.Info("unlink command called")

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

	fileParanoidPath := getParanoidPath(paranoidDirectory, filePath)
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
	inodeBytes, code, err := getFileInode(fileParanoidPath)
	if code != returncodes.OK {
		return code, err
	}

	//checking if we have access to the file.
	err = syscall.Access(path.Join(paranoidDirectory, "contents", string(inodeBytes)), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		return returncodes.EACCES, errors.New("could not access " + filePath)
	}

	// removing filename
	Log.Verbose("unlink : deleting file " + fileParanoidPath)
	err = os.Remove(fileParanoidPath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error removing file in names:", err)
	}

	// getting inode contents
	inodePath := path.Join(paranoidDirectory, "inodes", string(inodeBytes))
	Log.Verbose("unlink : reading file " + inodePath)
	inodeContents, err := ioutil.ReadFile(inodePath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error reading inodes contents:", err)
	}

	inodeData := &inode{}
	Log.Verbose("unlink unmarshaling ", string(inodeContents))
	err = json.Unmarshal(inodeContents, &inodeData)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error unmarshaling json:", err)
	}

	if inodeData.Count == 1 {
		// remove inode and contents
		contentsPath := path.Join(paranoidDirectory, "contents", string(inodeBytes))
		Log.Verbose("unlink : removing file " + contentsPath)
		err = os.Remove(contentsPath)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error removing contents:", err)
		}

		Log.Verbose("unlink : removing file " + inodePath)
		err = os.Remove(inodePath)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error removing inode:", err)
		}
	} else {
		// subtracting one from inode count and saving
		inodeData.Count--
		Log.Verbose("unlink : truncating file " + inodePath)
		err = os.Truncate(inodePath, 0)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error truncating inode path:", err)
		}

		dataToWrite, err := json.Marshal(inodeData)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error marshalling json:", err)
		}

		Log.Verbose("unlink : writing to file " + inodePath)
		err = ioutil.WriteFile(inodePath, dataToWrite, 0777)
		if err != nil {
			return returncodes.EUNEXPECTED, fmt.Errorf("error writing to inode file:", err)
		}
	}

	return returncodes.OK, nil
}
