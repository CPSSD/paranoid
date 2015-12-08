package commands

import (
	"encoding/json"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"os"
	"path"
	"syscall"
	"time"
)

type statInfo struct {
	Length int64       `json:"length",omitempty`
	Ctime  time.Time   `json:"ctime",omitempty`
	Mtime  time.Time   `json:"mtime",omitempty`
	Atime  time.Time   `json:"atime",omitempty`
	Mode   os.FileMode `json:"mode",omitempty`
}

//StatCommand prints a json object containing information on the file given as args[1] in pfs directory args[0] to Stdout
//Includes the length of the file, ctime, mtime and atime
func StatCommand(args []string) {
	Log.Info("Stat command called")
	if len(args) < 2 {
		Log.Fatal("Not enough arguments")
	}
	directory := args[0]
	Log.Verbose("stat : given directory", directory)

	getFileSystemLock(directory, sharedLock)
	defer unLockFileSystem(directory)

	namepath := getParanoidPath(directory, args[1])
	namePathType := getFileType(namepath)
	if namePathType == typeENOENT {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}

	fileNameBytes, code := getFileInode(namepath)
	if code != returncodes.OK {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(code))
		return
	}
	fileName := string(fileNameBytes)
	contentsFile := path.Join(directory, "contents", fileName)

	fi, err := os.Lstat(contentsFile)
	if err != nil {
		Log.Fatal("error Lstating file:", err)
	}

	stat := fi.Sys().(*syscall.Stat_t)
	atime := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
	ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
	var mode os.FileMode
	switch namePathType {
	case typeFile:
		mode = os.FileMode(stat.Mode)
	case typeSymlink:
		mode = os.FileMode(fi.Mode())
	default:
		mode = os.FileMode(syscall.S_IFDIR | fi.Mode().Perm())
	}

	statData := &statInfo{
		Length: fi.Size(),
		Mtime:  fi.ModTime(),
		Ctime:  ctime,
		Atime:  atime,
		Mode:   mode}

	jsonData, err := json.Marshal(statData)
	if err != nil {
		Log.Fatal("error marshalling statData:", err)
	}
	Log.Verbose("stat : returning", string(jsonData))
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
	io.WriteString(os.Stdout, string(jsonData))
}
