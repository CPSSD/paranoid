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

	if !checkFileExists(path.Join(directory, "names", args[1])) {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}

	fileNameBytes, err := ioutil.ReadFile(path.Join(directory, "names", args[1]))
	checkErr("access", err)
	fileName := string(fileNameBytes)

	getFileLock(directory, fileName, exclusiveLock)
	defer unLockFile(directory, fileName)

	mode, err := strconv.Atoi(args[2])
	checkErr("access", err)
	err = syscall.Access(path.Join(directory, "contents", fileName), uint32(mode))

	if err != nil {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EACCES))
		return
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
