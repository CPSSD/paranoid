package commands

import (
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"os"
	"path"
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

	linkInode, code := getFileInode(link)
	if code != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(code))
		return
	}

	inodePath := path.Join(directory, "inodes", string(linkInode))

	linkOriginBytes, err := ioutil.ReadFile(inodePath)
	checkErr("readlink", err)

	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
	io.WriteString(os.Stdout, string(linkOriginBytes))
}
