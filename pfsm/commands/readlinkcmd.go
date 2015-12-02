package commands

import (
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"os"
)

// ReadlinkCommand reads the value of the symbolic link
// args[0] is the init point and args[1] is the link
func ReadlinkCommand(args []string) {
	Log.Verbose("readlink called")
	if len(args) < 2 {
		Log.Fatal("not enough arguments")
	}

	directory := args[0]
	link := getParanoidPath(directory, args[1])

	Log.Verbose("readlink: given directory", directory)

	getFileSystemLock(directory, sharedLock)
	defer unLockFileSystem(directory)

	linkType := getFileType(link)
	if linkType == typeENOENT {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
	}

	// Not actually reading the inode, just getting
	linkOriginBytes, code := getFileInode(link)
	if code != returncodes.OK {
		io.WriteString(os.Stderr, returncodes.GetReturnCode(code))
		return
	}

	io.WriteString(os.Stderr, string(linkOriginBytes))
}
