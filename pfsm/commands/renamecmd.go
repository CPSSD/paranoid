package commands

import (
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"log"
	"os"
	"path"
	"syscall"
)

// RenameCommand is called when renaming a file
func RenameCommand(args []string) {
	verboseLog("rename command called")
	if len(args) < 3 {
		log.Fatalln("Not enough arguments!")
	}

	directory := args[0]
	oldFilePath := getParanoidPath(directory, args[1])
	newFilePath := getParanoidPath(directory, args[2])

	getFileSystemLock(directory, exclusiveLock)
	defer unLockFileSystem(directory)

	if !checkFileExists(oldFilePath) {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}

	if _, err := os.Stat(newFilePath); !os.IsNotExist(err) {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EEXIST))
		return
	}

	//Check if we have access to the file to be renamed
	if isDirectory(oldFilePath) {
		err := os.Rename(oldFilePath, newFilePath)
		checkErr("rename", err)
		oldInfoFilePath := path.Join(newFilePath, (path.Base(oldFilePath) + "-info"))
		newInfoFilePath := path.Join(newFilePath, (path.Base(newFilePath) + "-info"))
		err = os.Rename(oldInfoFilePath, newInfoFilePath)
		checkErr("rename", err)
	} else {
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
	}

	if !Flags.Network {
		sendToServer(directory, "rename", args[1:], nil)
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
