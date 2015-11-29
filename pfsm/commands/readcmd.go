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

//ReadCommand reads data from a file given as args[1] in pfs directory args[0] and prints it to Stdout
//Can also be given an offset and length as args[2] and args[3] otherwise it reads the whole file
func ReadCommand(args []string) {
	verboseLog("read command called")
	if len(args) < 2 {
		log.Fatalln("Not enough arguments!")
	}

	directory := args[0]
	verboseLog("read : given directory = " + directory)

	namepath := getParanoidPath(directory, args[1])

	getFileSystemLock(directory, sharedLock)
	defer unLockFileSystem(directory)

	if !checkFileExists(namepath) {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}

	fileNameBytes, err := ioutil.ReadFile(namepath)
	checkErr("read", err)
	fileName := string(fileNameBytes)

	err = syscall.Access(path.Join(directory, "contents", fileName), getAccessMode(syscall.O_RDONLY))
	if err != nil {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EACCES))
		return
	}

	getFileLock(directory, fileName, sharedLock)
	defer unLockFile(directory, fileName)

	file, err := os.OpenFile(path.Join(directory, "contents", fileName), os.O_RDONLY, 0777)
	checkErr("read", err)
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))

	if len(args) == 2 {
		verboseLog("read : reading whole file")
		bytesRead := make([]byte, 1024)
		for {
			n, err := file.Read(bytesRead)
			if err != io.EOF {
				checkErr("read", err)
			}
			io.WriteString(os.Stdout, string(bytesRead[:n]))
			if n < 1024 {
				break
			}
		}
	} else if len(args) > 2 {
		bytesRead := make([]byte, 1024)
		maxRead := 100000000
		if len(args) > 3 {
			verboseLog("read : " + args[3] + " bytes starting at " + args[2])
			maxRead, err = strconv.Atoi(args[3])
			checkErr("read", err)
		} else {
			verboseLog("read : from " + args[2] + " to end of file")
		}

		off, err := strconv.Atoi(args[2])
		checkErr("read", err)
		offset := int64(off)
		for {
			n, err := file.ReadAt(bytesRead, offset)
			if n > maxRead {
				bytesRead = bytesRead[0:maxRead]
				io.WriteString(os.Stdout, string(bytesRead))
				break
			}

			maxRead = maxRead - n
			if err == io.EOF {
				io.WriteString(os.Stdout, string(bytesRead[:n]))
				break
			}
			checkErr("read", err)
			io.WriteString(os.Stdout, string(bytesRead[:n]))
		}
	}
}
