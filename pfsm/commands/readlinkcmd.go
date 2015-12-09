package commands

import (
	"encoding/json"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// ReadlinkCommand reads the value of the symbolic link
// args[0] is the init point and args[1] is the link
func ReadlinkCommand(args []string) {
	Log.Info("readlink called")
	if len(args) < 2 {
		Log.Fatal("not enough arguments")
	}

	directory := args[0]
	link := getParanoidPath(directory, args[1])
	fileType := getFileType(directory, link)

	if fileType == typeENOENT {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}

	if fileType == typeDir {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EISDIR))
		return
	}

	if fileType == typeFile {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EIO))
		return
	}

	Log.Verbose("readlink: given directory", directory)

	getFileSystemLock(directory, sharedLock)
	defer unLockFileSystem(directory)

	linkInode, code := getFileInode(link)
	if code != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(code))
		return
	}

	inodePath := path.Join(directory, "inodes", string(linkInode))

	inodeContents, err := ioutil.ReadFile(inodePath)
	if err != nil {
		Log.Fatal("error reading link:", err)
	}

	inodeData := &inode{}
	Log.Verbose("unlink unmarshaling ", string(inodeContents))
	err = json.Unmarshal(inodeContents, &inodeData)
	if err != nil {
		Log.Fatal("error unmarshaling json ", err)
	}

	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
	io.WriteString(os.Stdout, string(inodeData.Link))
}
