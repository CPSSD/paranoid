package commands

import (
	"encoding/json"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"syscall"
)

// UnlinkCommand removes a filename link from an inode.
// if that is the only remaining link to the inode it removes the inode and its contents
func UnlinkCommand(args []string) {
	verboseLog("unlink command called")
	if len(args) < 2 {
		log.Fatalln("Not enough arguments!")
	}

	directory := args[0]
	fileNamePath := getParanoidPath(directory, args[1])

	verboseLog("unlink : directory given = " + directory)

	getFileSystemLock(directory, exclusiveLock)
	defer unLockFileSystem(directory)

	// checking if file exists
	if !checkFileExists(fileNamePath) {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
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
	verboseLog("unlink : deleting file " + fileNamePath)
	err = os.Remove(fileNamePath)
	checkErr("unlink", err)

	// getting inode contents
	inodePath := path.Join(directory, "inodes", string(inodeBytes))
	verboseLog("unlink : reading file " + inodePath)
	inodeContents, err := ioutil.ReadFile(inodePath)
	checkErr("unlink", err)
	inodeData := &inode{}
	err = json.Unmarshal(inodeContents, &inodeData)
	checkErr("unlink", err)

	if inodeData.Count == 1 {
		// remove inode and contents
		contentsPath := path.Join(directory, "contents", string(inodeBytes))
		verboseLog("unlink : removing file " + contentsPath)
		err = os.Remove(contentsPath)
		checkErr("unlink", err)
		verboseLog("unlink : removing file " + inodePath)
		err = os.Remove(inodePath)
		checkErr("unlink", err)
	} else {
		// subtracting one from inode count and saving
		inodeData.Count--
		verboseLog("unlink : truncating file " + inodePath)
		err = os.Truncate(inodePath, 0)
		checkErr("unlink", err)
		dataToWrite, err := json.Marshal(inodeData)
		checkErr("unlink", err)
		verboseLog("unlink : writing to file " + inodePath)
		err = ioutil.WriteFile(inodePath, dataToWrite, 0777)
		checkErr("unlink", err)
	}

	if !Flags.Network {
		sendToServer(directory, "unlink", args[1:], nil)
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
