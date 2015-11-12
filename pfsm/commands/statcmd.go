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
	"time"
)

type statInfo struct {
	Length int64       `json:"length",omitempty`
	Ctime  time.Time   `json:"ctime",omitempty`
	Mtime  time.Time   `json:"mtime",omitempty`
	Atime  time.Time   `json:"atime",omitempty`
	Perms  os.FileMode `json:"perms",omitempty`
}

//StatCommand prints a json object containing information on the file given as args[1] in pfs directory args[0] to Stdout
//Includes the length of the file, ctime, mtime and atime
func StatCommand(args []string) {
	verboseLog("stat command called")
	if len(args) < 2 {
		log.Fatalln("Not enough arguments!")
	}
	directory := args[0]
	verboseLog("stat : given directory = " + directory)

	getFileSystemLock(directory, sharedLock)
	defer unLockFileSystem(directory)

	if !checkFileExists(path.Join(directory, "names", args[1])) {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}

	fileNameBytes, err := ioutil.ReadFile(path.Join(directory, "names", args[1]))
	checkErr("stat", err)
	fileName := string(fileNameBytes)

	getFileLock(directory, fileName, sharedLock)
	defer unLockFile(directory, fileName)

	fi, err := os.Stat(path.Join(directory, "contents", fileName))
	checkErr("stat", err)

	stat := fi.Sys().(*syscall.Stat_t)
	atime := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
	ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
	statData := &statInfo{
		Length: fi.Size(),
		Mtime:  fi.ModTime(),
		Ctime:  ctime,
		Atime:  atime,
		Perms:  fi.Mode().Perm()}

	jsonData, err := json.Marshal(statData)
	checkErr("stat", err)
	verboseLog("stat : returning " + string(jsonData))
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
	io.WriteString(os.Stdout, string(jsonData))
}
