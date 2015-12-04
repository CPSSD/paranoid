package commands

import (
	"encoding/json"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

type inode struct {
	Count int    `json:"count"`
	Inode string `json:"inode"`
}

//CreatCommand creates a new file with the name args[1] in the pfs directory args[0]
func CreatCommand(args []string) {
	Log.Verbose("creat command called")
	if len(args) < 3 {
		log.Fatalln("Not enough arguments!")
	}
	directory := args[0]
	Log.Verbose("creat : directory = " + directory)

	getFileSystemLock(directory, exclusiveLock)
	defer unLockFileSystem(directory)

	namepath := getParanoidPath(directory, args[1])

	if getFileType(namepath) != typeENOENT {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EEXIST))
		return
	}
	Log.Verbose("creat : creating file " + args[1])

	uuidbytes := generateNewInode()
	uuidstring := string(uuidbytes)
	Log.Verbose("creat : uuid = " + uuidstring)

	perms, err := strconv.ParseInt(args[2], 8, 32)
	checkErr("creat", err)
	err = ioutil.WriteFile(namepath, uuidbytes, 0600)
	checkErr("creat", err)

	nodeData := &inode{
		Inode: uuidstring,
		Count: 1}
	jsonData, err := json.Marshal(nodeData)
	checkErr("creat", err)
	err = ioutil.WriteFile(path.Join(directory, "inodes", uuidstring), jsonData, 0600)
	checkErr("creat", err)

	contentsFile, err := os.Create(path.Join(directory, "contents", uuidstring))
	contentsFile.Chmod(os.FileMode(perms))
	checkErr("creat", err)

	if !Flags.Network {
		sendToServer(directory, "creat", args[1:], nil)
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
