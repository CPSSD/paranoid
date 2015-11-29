package commands

import (
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"log"
	"os"
	"path"
)

// RmdirCommand removes a directory
func RmdirCommand(args []string) {
	verboseLog("rmdir command called")
	if len(args) < 2 {
		log.Fatalln("Not enough arguments!")
	}

	directory := args[0]
	getFileSystemLock(directory, exclusiveLock)
	defer unLockFileSystem(directory)

	dirName, dirToDelete := getParanoidPath(directory, args[1])
	dirInfoFilePath := path.Join(dirToDelete, (dirName + "-info"))
	_, dirInodeString, retCode := getFileInode(dirInfoFilePath)
	if retCode != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(retCode))
		return
	}
	dirInodePath := path.Join(directory, "inodes", dirInodeString)

	code := deleteFile(dirInfoFilePath)
	if code != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(retCode))
		return
	}
	code = deleteFile(dirInodePath)
	if code != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(retCode))
		return
	}
	code = deleteFile(dirToDelete)
	if code != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(retCode))
		return
	}

	if !Flags.Network {
		sendToServer(directory, "rmdir", args[1:], nil)
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
