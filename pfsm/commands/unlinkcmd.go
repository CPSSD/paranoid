package commands

import (
	"encoding/json"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"syscall"
)

// UnlinkCommand removes a filename link from an inode.
// if that is the only remaining link to the inode it removes the inode and its contents
func UnlinkCommand(args []string) {
	Log.Info("unlink command called")
	if len(args) < 2 {
		Log.Fatal("Not enough arguments!")
	}

	directory := args[0]
	fileNamePath := getParanoidPath(directory, args[1])
	fileNamePathType := getFileType(fileNamePath)

	Log.Verbose("unlink : directory given = " + directory)

	getFileSystemLock(directory, exclusiveLock)
	defer unLockFileSystem(directory)

	// checking if file exists
	if fileNamePathType == typeENOENT {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}
	if fileNamePathType == typeDir {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EISDIR))
		return
	}

	// getting file inode
	inodeBytes, code := getFileInode(fileNamePath)
	if code != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(code))
		return
	}

	//checking if we have access to the file.
	err := syscall.Access(path.Join(directory, "contents", string(inodeBytes)), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EACCES))
		return
	}

	// removing filename
	Log.Verbose("unlink : deleting file " + fileNamePath)
	err = os.Remove(fileNamePath)
	if err != nil {
		Log.Fatal("error removing file in names:", err)
	}

	// getting inode contents
	inodePath := path.Join(directory, "inodes", string(inodeBytes))
	Log.Verbose("unlink : reading file " + inodePath)
	inodeContents, err := ioutil.ReadFile(inodePath)
	if err != nil {
		Log.Fatal("error reading inodes contents ", inodeContents)
	}

	inodeData := &inode{}
	Log.Verbose("unlink unmarshaling ", string(inodeContents))
	err = json.Unmarshal(inodeContents, &inodeData)
	if err != nil {
		Log.Fatal("error unmarshaling json ", err)
	}

	if inodeData.Count == 1 {
		// remove inode and contents
		contentsPath := path.Join(directory, "contents", string(inodeBytes))
		Log.Verbose("unlink : removing file " + contentsPath)
		err = os.Remove(contentsPath)
		if err != nil {
			Log.Fatal("error removing contents:", err)
		}

		Log.Verbose("unlink : removing file " + inodePath)
		err = os.Remove(inodePath)
		if err != nil {
			Log.Fatal("error removing inode:", err)
		}
	} else {
		// subtracting one from inode count and saving
		inodeData.Count--
		Log.Verbose("unlink : truncating file " + inodePath)
		err = os.Truncate(inodePath, 0)
		if err != nil {
			Log.Fatal("error truncating inode path:", err)
		}

		dataToWrite, err := json.Marshal(inodeData)
		if err != nil {
			Log.Fatal("error marshalling json:", err)
		}
		Log.Verbose("unlink : writing to file " + inodePath)
		err = ioutil.WriteFile(inodePath, dataToWrite, 0777)
		if err != nil {
			Log.Fatal("error writing to inode file:", err)
		}
	}

	if !Flags.Network {
		sendToServer(directory, "unlink", args[1:], nil)
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
