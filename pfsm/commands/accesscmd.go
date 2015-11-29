package commands

import (
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"syscall"
)

//AccessCommand is used by fuse to check if it has access to a given file.
func AccessCommand(args []string) {
	verboseLog("access command given")
	if len(args) < 3 {
		log.Fatalln("Not enough arguments!")
	}
	directory := args[0]
	verboseLog("access : given directory = " + directory)

	getFileSystemLock(directory, sharedLock)
	defer unLockFileSystem(directory)

	_, namePath := getParanoidPath(directory, args[1])

	if !checkFileExists(namePath) {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}

	fileNameBytes, err := ioutil.ReadFile(namePath)
	checkErr("access", err)
	fileName := string(fileNameBytes)

	mode, err := strconv.Atoi(args[2])
	checkErr("access", err)
	err = syscall.Access(path.Join(directory, "contents", fileName), uint32(mode))

	if err != nil {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EACCES))
		return
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
