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

//ChmodCommand is used to change the permissions of a file.
func ChmodCommand(args []string) {
	verboseLog("chmod command given")
	if len(args) < 3 {
		log.Fatalln("Not enough arguments!")
	}
	directory := args[0]
	verboseLog("chmod : given directory = " + directory)

	getFileSystemLock(directory, exclusiveLock)
	defer unLockFileSystem(directory)

	if !checkFileExists(path.Join(directory, "names", args[1])) {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}

	fileNameBytes, err := ioutil.ReadFile(path.Join(directory, "names", args[1]))
	checkErr("chmod", err)
	fileName := string(fileNameBytes)

	err = syscall.Access(path.Join(directory, "contents", fileName), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EACCES))
		return
	}

	verboseLog("chmod : changing permissions of " + fileName + " to " + args[2])
	perms, err := strconv.ParseInt(args[2], 8, 32)
	checkErr("chmod", err)

	contentsFile, err := os.OpenFile(path.Join(directory, "contents", fileName), os.O_WRONLY, 0777)
	checkErr("chmod", err)
	err = contentsFile.Chmod(os.FileMode(perms))
	checkErr("chmod", err)

	if !Flags.Network {
		sendToServer(directory, "chmod", args[1:], nil)
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
