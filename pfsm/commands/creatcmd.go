package commands

import (
	"encoding/json"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

//CreatCommand creates a new file with the name args[1] in the pfs directory args[0]
func CreatCommand(args []string) {
	Log.Info("creat command called")
	if len(args) < 3 {
		Log.Fatal("Not enough arguments!")
	}
	directory := args[0]
	Log.Verbose("creat : directory = " + directory)

	getFileSystemLock(directory, exclusiveLock)
	defer unLockFileSystem(directory)

	namepath := getParanoidPath(directory, args[1])

	if getFileType(directory, namepath) != typeENOENT {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EEXIST))
		return
	}
	Log.Verbose("creat : creating file " + args[1])

	uuidbytes := generateNewInode()
	uuidstring := string(uuidbytes)
	Log.Verbose("creat : uuid = " + uuidstring)

	perms, err := strconv.ParseInt(args[2], 8, 32)
	if err != nil {
		Log.Fatal("error converting mode from string to int:", err)
	}

	err = ioutil.WriteFile(namepath, uuidbytes, 0600)
	if err != nil {
		Log.Fatal("error writing name file", err)
	}

	nodeData := &inode{
		Inode: uuidstring,
		Count: 1}
	jsonData, err := json.Marshal(nodeData)
	if err != nil {
		Log.Fatal("error marshalling json:", err)
	}

	err = ioutil.WriteFile(path.Join(directory, "inodes", uuidstring), jsonData, 0600)
	if err != nil {
		Log.Fatal("error writing inodes file:", err)
	}

	contentsFile, err := os.Create(path.Join(directory, "contents", uuidstring))
	if err != nil {
		Log.Fatal("error creating contents file:", err)
	}

	err = contentsFile.Chmod(os.FileMode(perms))
	if err != nil {
		Log.Fatal("error changing file permissions:", err)
	}

	if !Flags.Network {
		sendToServer(directory, "creat", args[1:], nil)
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
