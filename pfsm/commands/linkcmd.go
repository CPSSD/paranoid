package commands

import (
	"encoding/json"
	"errors"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io/ioutil"
	"os"
	"path"
)

// LinkCommand creates a link of a file.
// args[0] is the init point, args[1] is the existing file name, args[2] is the target file name
func LinkCommand(directory, existingFileName, targetFileName string) (returnCode int, returnError error) {
	Log.Info("link command called")

	existingFilePath := getParanoidPath(directory, existingFileName)
	targetFilePath := getParanoidPath(directory, targetFileName)

	Log.Verbose("link : given directory = " + directory)

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

	existingFileType := getFileType(directory, existingFilePath)
	if existingFileType == typeENOENT {
		return returncodes.ENOENT, errors.New("existing file does not exist")
	}

	if existingFileType == typeDir {
		return returncodes.EISDIR, errors.New("existing file is a directory")
	}

	if existingFileType == typeSymlink {
		return returncodes.EIO, errors.New("existing file is a symlink")
	}

	if getFileType(directory, targetFilePath) != typeENOENT {
		return returncodes.EEXIST, errors.New("target file already exists")
	}

	// getting inode and fileMode of existing file
	inodeBytes, code, err := getFileInode(existingFilePath)
	if code != returncodes.OK {
		return code, err
	}

	fileInfo, err := os.Stat(existingFilePath)
	if err != nil {
		return returncodes.EUNEXPECTED, errors.New("error stating existing file:", err)
	}
	fileMode := fileInfo.Mode()

	// creating target file pointing to same inode
	err = ioutil.WriteFile(targetFilePath, inodeBytes, fileMode)
	if err != nil {
		return returncodes.EUNEXPECTED, errors.New("error writing to names file:", err)
	}

	// getting contents of inode
	inodePath := path.Join(directory, "inodes", string(inodeBytes))
	Log.Verbose("link : reading file " + inodePath)
	inodeContents, err := ioutil.ReadFile(inodePath)
	if err != nil {
		return returncodes.EUNEXPECTED, errors.New("error reading inode:", err)
	}

	nodeData := &inode{}
	err = json.Unmarshal(inodeContents, &nodeData)
	if err != nil {
		return returncodes.EUNEXPECTED, errors.New("error unmarshalling inode data:", err)
	}

	// itterating count and saving
	nodeData.Count++
	Log.Verbose("link : opening file " + inodePath)
	openedFile, err := os.OpenFile(inodePath, os.O_WRONLY, 0600)
	if err != nil {
		return returncodes.EUNEXPECTED, errors.New("error opening file:", err)
	}

	Log.Verbose("link : truncating file " + inodePath)
	err = openedFile.Truncate(0)
	if err != nil {
		return returncodes.EUNEXPECTED, errors.New("error truncating file:", err)
	}

	newJSONData, err := json.Marshal(&nodeData)
	if err != nil {
		return returncodes.EUNEXPECTED, errors.New("error marshalling json:", err)
	}

	Log.Verbose("link : writing to file " + inodePath)
	_, err = openedFile.Write(newJSONData)
	if err != nil {
		return returncodes.EUNEXPECTED, errors.New("error writing to inode file:", err)
	}

	// closing file
	Log.Verbose("link : closing file " + inodePath)
	err = openedFile.Close()
	if err != nil {
		return returncodes.EUNEXPECTED, errors.New("error closing file:", err)
	}

	if !Flags.Network {
		sendToServer(directory, "link", args[1:], nil)
	}
	return returncodes.OK, nil
}
