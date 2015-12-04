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
	Log.Verbose("mkdir command called")
	if len(args) < 3 {
		log.Fatalln("Not enough arguments!")
	}

	directory := args[0]
	getFileSystemLock(directory, exclusiveLock)
	defer unLockFileSystem(directory)

	dirPath := getParanoidPath(directory, args[1])
	dirInfoPath := path.Join(dirPath, "info")
	inodeBytes := generateNewInode()
	inodeString := string(inodeBytes)
	inodePath := path.Join(directory, "inodes", inodeString)
	contentsPath := path.Join(directory, "contents", inodeString)
	mode, err := strconv.ParseInt(args[2], 8, 32)
	checkErr("mkdir", err)

	if getFileType(dirPath) != typeENOENT {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EEXIST))
		return
	}
	err = os.Mkdir(dirPath, os.FileMode(mode))
	checkErr("mkdir", err)

	contentsFile, err := os.Create(contentsPath)
	checkErr("mkdir", err)
	err = contentsFile.Chmod(os.FileMode(mode))
	checkErr("mkdir", err)

	dirInfoFile, err := os.Create(dirInfoPath)
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
