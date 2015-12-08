package commands

import (
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"syscall"
)

//ChmodCommand is used to change the permissions of a file.
func ChmodCommand(args []string) {
	Log.Info("chmod command given")
	if len(args) < 3 {
		log.Fatalln("Not enough arguments!")
	}
	directory := args[0]
	Log.Verbose("chmod : given directory = " + directory)

	getFileSystemLock(directory, exclusiveLock)
	defer unLockFileSystem(directory)

	namepath := getParanoidPath(directory, args[1])

	fileType := getFileType(namepath)
	if fileType == typeENOENT {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}

	if fileType == typeSymlink {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EIO))
		return
	}

	fileNameBytes, code := getFileInode(namepath)
	if code != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(code))
		return
	}
	fileName := string(fileNameBytes)

	err := syscall.Access(path.Join(directory, "contents", fileName), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EACCES))
		return
	}

	Log.Verbose("chmod : changing permissions of " + fileName + " to " + args[2])
	perms, err := strconv.ParseInt(args[2], 8, 32)
	if err != nil {
		Log.Fatal("error converting permissions from string to int:", err)
	}

	contentsFile, err := os.OpenFile(path.Join(directory, "contents", fileName), os.O_WRONLY, 0777)
	if err != nil {
		Log.Fatal("error opening contents file:", err)
	}

	err = contentsFile.Chmod(os.FileMode(perms))
	if err != nil {
		Log.Fatal("error changing permissions:", err)
	}

	if !Flags.Network {
		sendToServer(directory, "chmod", args[1:], nil)
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
