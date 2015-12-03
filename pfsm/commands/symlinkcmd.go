package commands

import (
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"os"
)

// SymlinkCommand creates a symbolic link
// args[0] is the init point, args[1] is the existing file, args[2] is the target file
func SymlinkCommand(args []string) {
	Log.Info("symlink command called")
	if len(args) < 3 {
		Log.Fatal("not enough arguments")
	}

	directory := args[0]
	targetFilePath := getParanoidPath(directory, args[2])

	Log.Verbose("symlink: given directory ", directory)

	getFileSystemLock(directory, exclusiveLock)
	defer unLockFileSystem(directory)

	// Make sure the target file not existing, if it is, quit
	if getFileType(targetFilePath) != typeENOENT {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EEXIST))
		return
	}

	// Create a new file with content which is the
	// relative location of the existing file
	err := ioutil.WriteFile(targetFilePath, []byte(args[1]), 0777)
	checkErr("symlink", err)

	// Send to the server if not coming from the network
	if !Flags.Network {
		sendToServer(directory, "symlink", args[1:], nil)
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
