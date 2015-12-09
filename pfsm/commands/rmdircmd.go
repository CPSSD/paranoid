package commands

import (
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// RmdirCommand removes a directory
func RmdirCommand(args []string) {
	Log.Info("rmdir command called")
	if len(args) < 2 {
		Log.Fatal("Not enough arguments!")
	}

	directory := args[0]
	getFileSystemLock(directory, exclusiveLock)
	defer unLockFileSystem(directory)

	dirToDelete := getParanoidPath(directory, args[1])
	dirType := getFileType(directory, dirToDelete)
	if dirType == typeENOENT {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	} else if dirType != typeDir {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOTDIR))
		return
	} else {
		files, err := ioutil.ReadDir(dirToDelete)
		if err != nil {
			if os.IsPermission(err) {
				io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EACCES))
				return
			}
			Log.Fatal("error reading directory:", err)
		}
		if len(files) > 1 {
			io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOTEMPTY))
			return
		}
	}

	infoFileToDelete := path.Join(dirToDelete, "info")
	inodeBytes, err := getFileInode(dirToDelete)
	inodeString := string(inodeBytes)
	if err != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(err))
		return
	}
	inodeFileToDelete := path.Join(directory, "inodes", inodeString)
	contentsFileToDelete := path.Join(directory, "contents", inodeString)

	code := deleteFile(contentsFileToDelete)
	if code != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(code))
		return
	}

	code = deleteFile(inodeFileToDelete)
	if code != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(code))
		return
	}

	code = deleteFile(infoFileToDelete)
	if code != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(code))
		return
	}

	code = deleteFile(dirToDelete)
	if code != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(code))
		return
	}

	if !Flags.Network {
		sendToServer(directory, "rmdir", args[1:], nil)
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
