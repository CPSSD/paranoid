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

// LinkCommand creates a link of a file.
func LinkCommand(paranoidDirectory, existingFilePath, targetFilePath string) (returnCode int, returnError error) {
	Log.Info("link command called")
	existingParanoidPath := getParanoidPath(paranoidDirectory, existingFilePath)
	targetParanoidPath := getParanoidPath(paranoidDirectory, targetFilePath)

	Log.Verbose("link : given paranoidDirectory = " + paranoidDirectory)

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

	existingFileType, err := getFileType(paranoidDirectory, existingParanoidPath)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	if existingFileType == typeENOENT {
		return returncodes.ENOENT, errors.New("existing file " + existingFilePath + " does not exist")
	}

	if existingFileType == typeDir {
		return returncodes.EISDIR, errors.New("existing file " + existingFilePath + " is a paranoidDirectory")
	}

	if existingFileType == typeSymlink {
		return returncodes.EIO, errors.New("existing file " + existingFilePath + " is a symlink")
	}

	targetFileType, err := getFileType(paranoidDirectory, targetParanoidPath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error getting target file "+targetFilePath+" file type:", err)
	}

	if targetFileType != typeENOENT {
		return returncodes.EEXIST, errors.New("target file " + targetFilePath + " already exists")
	}

	// getting inode and fileMode of existing file
	inodeBytes, code, err := getFileInode(existingParanoidPath)
	if code != returncodes.OK {
		return code, err
	}

	fileInfo, err := os.Stat(existingParanoidPath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error stating existing file "+existingFilePath+":", err)
	}
	fileMode := fileInfo.Mode()

	// creating target file pointing to same inode
	err = ioutil.WriteFile(targetParanoidPath, inodeBytes, fileMode)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing to names file:", err)
	}

	// getting contents of inode
	inodePath := path.Join(paranoidDirectory, "inodes", string(inodeBytes))
	Log.Verbose("link : reading file " + inodePath)
	inodeContents, err := ioutil.ReadFile(inodePath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error reading inode:", err)
	}

	nodeData := &inode{}
	err = json.Unmarshal(inodeContents, &nodeData)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error unmarshalling inode data:", err)
	}

	// itterating count and saving
	nodeData.Count++
	Log.Verbose("link : opening file " + inodePath)
	openedFile, err := os.OpenFile(inodePath, os.O_WRONLY, 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error opening file:", err)
	}
	defer openedFile.Close()

	Log.Verbose("link : truncating file " + inodePath)
	err = openedFile.Truncate(0)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error truncating file:", err)
	}

	newJSONData, err := json.Marshal(&nodeData)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error marshalling json:", err)
	}

	Log.Verbose("link : writing to file " + inodePath)
	_, err = openedFile.Write(newJSONData)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error writing to inode file:", err)
	}

	// closing file
	Log.Verbose("link : closing file " + inodePath)
	err = openedFile.Close()
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error closing file:", err)
	}

	return returncodes.OK, nil
}
