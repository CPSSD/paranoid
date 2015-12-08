package commands

import (
	"encoding/json"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"os"
	"path"
	"strconv"
)

// MkdirCommand is called when making a directory
func MkdirCommand(args []string) {
	Log.Info("mkdir command called")
	if len(args) < 3 {
		Log.Fatal("Not enough arguments!")
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
	if err != nil {
		Log.Fatal("error converting mode from string to int:", err)
	}

	if getFileType(dirPath) != typeENOENT {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EEXIST))
		return
	}
	err = os.Mkdir(dirPath, os.FileMode(mode))
	if err != nil {
		Log.Fatal("error making directory:", err)
	}

	contentsFile, err := os.Create(contentsPath)
	if err != nil {
		Log.Fatal("error creating contents file:", err)
	}

	err = contentsFile.Chmod(os.FileMode(mode))
	if err != nil {
		Log.Fatal("error changing file permissions:", err)
	}

	dirInfoFile, err := os.Create(dirInfoPath)
	if err != nil {
		Log.Fatal("error creating info file:", err)
	}

	_, err = dirInfoFile.Write(inodeBytes)
	if err != nil {
		Log.Fatal("error writing to info file:", err)
	}

	inodeFile, err := os.Create(inodePath)
	if err != nil {
		Log.Fatal("error creating inode file:", err)
	}

	nodeData := &inode{
		Inode: inodeString,
		Count: 1}
	jsonData, err := json.Marshal(nodeData)
	if err != nil {
		Log.Fatal("error marshalling json:", err)
	}

	_, err = inodeFile.Write(jsonData)
	if err != nil {
		Log.Fatal("error writing to inode file", err)
	}

	if !Flags.Network {
		sendToServer(directory, "mkdir", args[1:], nil)
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
