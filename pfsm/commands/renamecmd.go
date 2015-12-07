package commands

import (
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"os"
	"path"
	"syscall"
)

// RenameCommand is called when renaming a file
func RenameCommand(args []string) {
	Log.Verbose("rename command called")
	if len(args) < 3 {
		Log.Fatal("Not enough arguments!")
	}

	directory := args[0]
	oldFilePath := getParanoidPath(directory, args[1])
	newFilePath := getParanoidPath(directory, args[2])
	oldFileType := getFileType(oldFilePath)
	newFileType := getFileType(newFilePath)

	getFileSystemLock(directory, exclusiveLock)
	defer unLockFileSystem(directory)

	if oldFileType == typeENOENT {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}
	if newFileType != typeENOENT {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EEXIST))
		return
	}

	inodeBytes, code := getFileInode(oldFilePath)
	if code != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(code))
		return
	}
	err := syscall.Access(path.Join(directory, "contents", string(inodeBytes)), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EACCES))
		return
	}

	err = os.Rename(oldFilePath, newFilePath)
	checkErr("rename", err)

	if !Flags.Network {
		sendToServer(directory, "rename", args[1:], nil)
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
