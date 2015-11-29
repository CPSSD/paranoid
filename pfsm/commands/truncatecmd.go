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

//TruncateCommand reduces the file given as args[1] in the paranoid-direcory args[0] to the size given in args[2]
func TruncateCommand(args []string) {
	verboseLog("truncate command given")
	if len(args) < 3 {
		log.Fatalln("Not enough arguments!")
	}
	directory := args[0]
	verboseLog("truncate : given directory = " + directory)

	getFileSystemLock(directory, sharedLock)
	defer unLockFileSystem(directory)

	namepath := getParanoidPath(directory, args[1])

	if !checkFileExists(namepath) {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}

	fileNameBytes, err := ioutil.ReadFile(namepath)
	checkErr("truncate", err)
	fileName := string(fileNameBytes)

	err = syscall.Access(path.Join(directory, "contents", fileName), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EACCES))
		return
	}

	getFileLock(directory, fileName, exclusiveLock)
	defer unLockFile(directory, fileName)

	verboseLog("truncate : truncating " + fileName)
	newsize, err := strconv.Atoi(args[2])
	checkErr("truncate", err)
	contentsFile, err := os.OpenFile(path.Join(directory, "contents", fileName), os.O_WRONLY, 0777)
	checkErr("truncate", err)
	err = contentsFile.Truncate(int64(newsize))
	checkErr("truncate", err)

	if !Flags.Network {
		sendToServer(directory, "truncate", args[1:], nil)
	}
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
