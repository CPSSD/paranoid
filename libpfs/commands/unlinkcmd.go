package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/pfsd/pnetclient"
	"io/ioutil"
	"os"
	"path"
	"syscall"
)

// UnlinkCommand removes a filename link from an inode.
func UnlinkCommand(directory, fileName string, sendOverNetwork bool) (returnCode int, returnError error) {
	Log.Info("unlink command called")

	err := getFileSystemLock(directory, exclusiveLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	defer func() {
		err := unLockFileSystem(directory)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
		}
	}()

	fileNamePath := getParanoidPath(directory, fileName)
	fileNamePathType, err := getFileType(directory, fileNamePath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	Log.Verbose("unlink : directory given = " + directory)

	// checking if file exists
	if fileNamePathType == typeENOENT {
		return returncodes.ENOENT, errors.New(fileName + " does not exist")
	}

	if fileNamePathType == typeDir {
		return returncodes.EISDIR, errors.New(fileName + " is a directory")
	}

	// getting file inode
	inodeBytes, code, err := getFileInode(fileNamePath)
	if code != returncodes.OK {
		return code, err
	}

	//checking if we have access to the file.
	err = syscall.Access(path.Join(directory, "contents", string(inodeBytes)), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		return returncodes.EACCES, errors.New("could not access " + fileName)
	}

	// removing filename
	Log.Verbose("unlink : deleting file " + fileNamePath)
	err = os.Remove(fileNamePath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error removing file in names:", err)
	}

	// getting inode contents
	inodePath := path.Join(directory, "inodes", string(inodeBytes))
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
		contentsPath := path.Join(directory, "contents", string(inodeBytes))
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

	if sendOverNetwork {
		pnetclient.Unlink(fileName)
	}
	return returncodes.OK, nil
}
