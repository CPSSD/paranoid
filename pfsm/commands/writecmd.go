package commands

import (
	"github.com/cpssd/paranoid/pfsm/network"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"syscall"
)

//WriteCommand writes data from Stdin to the file given in args[1] in the pfs directory args[0]
//Can also be given an offset and length as args[2] and args[3] otherwise it writes from the start of the file
func WriteCommand(args []string) {
	verboseLog("write command given")
	if len(args) < 2 {
		log.Fatalln("Not enough arguments!")
	}
	directory := args[0]
	verboseLog("write : given directory = " + directory)

	getFileSystemLock(directory, sharedLock)
	defer unLockFileSystem(directory)

	if !checkFileExists(path.Join(directory, "names", args[1])) {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}

	fileNameBytes, err := ioutil.ReadFile(path.Join(directory, "names", args[1]))
	checkErr("write", err)
	fileName := string(fileNameBytes)

	getFileLock(directory, fileName, exclusiveLock)
	defer unLockFile(directory, fileName)

	err = syscall.Access(path.Join(directory, "contents", fileName), getAccessMode(syscall.O_WRONLY))
	if err != nil {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EACCES))
		return
	}

	verboseLog("write : wrting to " + fileName)
	fileData, err := ioutil.ReadAll(os.Stdin)
	checkErr("write", err)

	if len(args) == 2 {
		err = ioutil.WriteFile(path.Join(directory, "contents", fileName), fileData, 0777)
		checkErr("write", err)
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
		io.WriteString(os.Stdout, strconv.Itoa(len(fileData)))

		if !Flags.Network {
			network.Write(directory, args[1], nil, nil, string(fileData))
		}
	} else {
		contentsFile, err := os.OpenFile(path.Join(directory, "contents", fileName), os.O_WRONLY, 0777)
		checkErr("write", err)
		offset, err := strconv.Atoi(args[2])
		checkErr("write", err)

		if len(args) == 3 {
			err = contentsFile.Truncate(int64(offset))
			checkErr("write", err)
		} else {
			length, err := strconv.Atoi(args[3])
			checkErr("write", err)
			if len(fileData) > length {
				fileData = fileData[:length]
			} else if len(fileData) < length {
				emptyBytes := make([]byte, length-len(fileData))
				fileData = append(fileData, emptyBytes...)
			}
		}

		wroteLen, err := contentsFile.WriteAt(fileData, int64(offset))
		checkErr("write", err)
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
		io.WriteString(os.Stdout, strconv.Itoa(wroteLen))

		if len(args) == 3 {
			if !Flags.Network {
				network.Write(directory, args[1], &offset, nil, string(fileData))
			}
		} else {
			if !Flags.Network {
				length, err := strconv.Atoi(args[3])
				checkErr("write", err)
				network.Write(directory, args[1], &offset, &length, string(fileData))
			}
		}
	}
}
