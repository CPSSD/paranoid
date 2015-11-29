package commands

import (
	"encoding/json"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// LinkCommand creates a link of a file.
// args[0] is the init point, args[1] is the existing file name, args[2] is the target file name
func LinkCommand(args []string) {
	verboseLog("link command called")
	if len(args) < 3 {
		log.Fatalln("Not enough arguments!")
	}

	directory := args[0]
	_, existingFilePath := getParanoidPath(directory, args[1])
	_, targetFilePath := getParanoidPath(directory, args[2])

	verboseLog("link : given directory = " + directory)

	getFileSystemLock(directory, exclusiveLock)
	defer unLockFileSystem(directory)

	if !checkFileExists(existingFilePath) {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}
	if checkFileExists(targetFilePath) {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EEXIST))
		return
	}

	// getting inode and fileMode of existing file
	inodeBytes, inodeString, code := getFileInode(existingFilePath)
	if code != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(code))
		return
	}
	fileInfo, err := os.Stat(existingFilePath)
	checkErr("link", err)
	fileMode := fileInfo.Mode()

	// creating target file pointing to same inode
	err = ioutil.WriteFile(targetFilePath, inodeBytes, fileMode)
	checkErr("link", err)

	// getting contents of inode
	inodePath := path.Join(directory, "inodes", inodeString)
	verboseLog("link : reading file " + inodePath)
	inodeContents, err := ioutil.ReadFile(inodePath)
	checkErr("link", err)
	nodeData := &inode{}
	err = json.Unmarshal(inodeContents, &nodeData)
	checkErr("link", err)

	// itterating count and saving
	nodeData.Count++
	verboseLog("link : opening file " + inodePath)
	openedFile, err := os.OpenFile(inodePath, os.O_WRONLY, 0600)
	checkErr("link", err)
	verboseLog("link : truncating file " + inodePath)
	err = openedFile.Truncate(0)
	checkErr("link", err)
	newJSONData, err := json.Marshal(&nodeData)
	checkErr("link", err)
	verboseLog("link : writing to file " + inodePath)
	_, err = openedFile.Write(newJSONData)
	checkErr("link", err)

	// closing file
	verboseLog("link : closing file " + inodePath)
	err = openedFile.Close()
	checkErr("link", err)

	if !Flags.Network {
		sendToServer(directory, "link", args[1:], nil)
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
