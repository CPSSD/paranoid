package commands

import (
	"encoding/json"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"log"
	"os"
	"path"
	"strconv"
)

// MkdirCommand is called when making a directory
func MkdirCommand(args []string) {
	verboseLog("mkdir command called")
	if len(args) < 3 {
		log.Fatalln("Not enough arguments!")
	}

	directory := args[0]
	getFileSystemLock(directory, exclusiveLock)
	defer unLockFileSystem(directory)

	dirPath := getParanoidPath(directory, args[1])
	dirInfoPath := path.Join(dirPath, (path.Base(dirPath) + "-info"))
	inodeBytes, inodeString := generateNewInode()
	inodePath := path.Join(directory, "inodes", inodeString)
	mode, err := strconv.Atoi(args[2])
	checkErr("mkdir", err)

	if checkFileExists(dirPath) {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EEXIST))
		return
	}
	err = os.Mkdir(dirPath, os.FileMode(mode))
	checkErr("mkdir", err)

	dirInfoFile, err := os.Create(dirInfoPath)
	checkErr("mkdir", err)
	err = dirInfoFile.Chmod(os.FileMode(mode))
	checkErr("mkdir", err)
	_, err = dirInfoFile.Write(inodeBytes)
	checkErr("mkdir", err)

	inodeFile, err := os.Create(inodePath)
	checkErr("mkdir", err)
	nodeData := &inode{
		Inode: inodeString,
		Count: 1}
	jsonData, err := json.Marshal(nodeData)
	checkErr("mkdir", err)
	_, err = inodeFile.Write(jsonData)
	checkErr("mkdir", err)

	if !Flags.Network {
		sendToServer(directory, "mkdir", args[1:], nil)
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
