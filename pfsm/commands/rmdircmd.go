package commands

import (
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
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

	dirToDelete := getParanoidPath(directory, args[1])
	dirInfoFilePath := path.Join(dirToDelete, (path.Base(dirToDelete) + "-info"))
	dirInodeBytes, retCode := getFileInode(dirInfoFilePath)
	if retCode != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(retCode))
		return
	}
	if !checkFileExists(dirToDelete) {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}
	names, err := ioutil.ReadDir(dirToDelete)
	checkErr("rmdir", err)
	if len(names) > 1 {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOTEMPTY))
		return
	}
	checkErr("rmdir", err)
	dirInodePath := path.Join(directory, "inodes", string(dirInodeBytes))

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
