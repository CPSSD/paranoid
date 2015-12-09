package commands

import (
	"encoding/json"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"syscall"
	"time"
)

type timeInfo struct {
	Atime *time.Time `json:"atime",omitempty`
	Mtime *time.Time `json:"mtime",omitempty`
}

//UtimesCommand updates the acess time and modified time of a file
func UtimesCommand(args []string) {
	Log.Info("utimes command called")
	if len(args) < 2 {
		Log.Fatal("Not enough arguments!")
	}
	directory := args[0]
	Log.Verbose("utimes : given directory = " + directory)

	getFileSystemLock(directory, sharedLock)
	defer unLockFileSystem(directory)

	namepath := getParanoidPath(directory, args[1])

	if getFileType(directory, namepath) == typeENOENT {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		Log.Fatal("error reading input:", err)
	}
	times := timeInfo{}
	err = json.Unmarshal(input, &times)
	if err != nil {
		Log.Fatal("error unmarshalling times:", err)
	}

	fileNameBytes, code := getFileInode(namepath)
	if code != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(code))
		return
	}
	fileName := string(fileNameBytes)

	err = syscall.Access(path.Join(directory, "contents", fileName), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EACCES))
		return
	}

	getFileLock(directory, fileName, exclusiveLock)
	defer unLockFile(directory, fileName)

	file, err := os.Open(path.Join(directory, "contents", fileName))
	if err != nil {
		Log.Fatal("error opening contents file:", err)
	}

	fi, err := file.Stat()
	if err != nil {
		Log.Fatal("error stating file:", err)
	}
	stat := fi.Sys().(*syscall.Stat_t)
	oldatime := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
	oldmtime := fi.ModTime()
	if times.Atime == nil && times.Mtime == nil {
		Log.Fatal("utimes : no times to update!")
	}

	if times.Atime == nil {
		os.Chtimes(path.Join(directory, "contents", fileName), oldatime, *times.Mtime)
	} else if times.Mtime == nil {
		os.Chtimes(path.Join(directory, "contents", fileName), *times.Atime, oldmtime)
	} else {
		os.Chtimes(path.Join(directory, "contents", fileName), *times.Atime, *times.Mtime)
	}

	if !Flags.Network {
		sendToServer(directory, "utimes", args[1:], input)
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
