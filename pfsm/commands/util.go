package commands

import (
	"github.com/cpssd/paranoid/pfsm/icclient"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"syscall"
)

//Check if a given file exists
func checkFileExists(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

func getAccessMode(flags uint32) uint32 {
	switch flags {
	case syscall.O_RDONLY:
		return 4
	case syscall.O_WRONLY:
		return 2
	case syscall.O_RDWR:
		return 6
	default:
		return 7
	}
}

//verboseLog logs a message if the verbose command line flag was set.
func verboseLog(message string) {
	if Flags.Verbose {
		log.Println(message)
	}
}

//checkErr stops the execution of the program if the given error is not nil.
//Specifies the command where the error occured as cmd
func checkErr(cmd string, err error) {
	if err != nil {
		log.Fatalln(cmd, " error occured: ", err)
	}
}

//Types of locks
const (
	sharedLock = iota
	exclusiveLock
)

func getFileSystemLock(paranoidDir string, lockType int) {
	lockPath := path.Join(paranoidDir, "meta", "lock")
	file, err := os.Open(lockPath)
	if err != nil {
		log.Fatalln("Could not get meta/lock file discriptor")
	}
	defer file.Close()
	if lockType == sharedLock {
		syscall.Flock(int(file.Fd()), syscall.LOCK_SH)
	} else if lockType == exclusiveLock {
		syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
	}
}

func getFileLock(paranoidDir, fileName string, lockType int) {
	lockPath := path.Join(paranoidDir, "contents", fileName)
	file, err := os.Open(lockPath)
	if err != nil {
		log.Fatalln("Could not get file discriptor for lock file : ", fileName)
	}
	defer file.Close()
	if lockType == sharedLock {
		syscall.Flock(int(file.Fd()), syscall.LOCK_SH)
	} else if lockType == exclusiveLock {
		syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
	}
}

func unLockFileSystem(paranoidDir string) {
	lockPath := path.Join(paranoidDir, "meta", "lock")
	file, err := os.Open(lockPath)
	if err != nil {
		log.Fatalln("Could not get meta/lock file discriptor for unlock")
	}
	defer file.Close()
	syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
}

func unLockFile(paranoidDir, fileName string) {
	lockPath := path.Join(paranoidDir, "contents", fileName)
	file, err := os.Open(lockPath)
	if err != nil {
		log.Fatalln("Could not get file discriptor for unlock file :", fileName)
	}
	defer file.Close()
	syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
}

func sendToServer(paranoidDir, command string, args []string, data []byte) {
	if data == nil {
		icclient.SendMessage(paranoidDir, command, args)
	} else {
		icclient.SendMessageWithData(paranoidDir, command, args, data)
	}
}

func getParanoidPath(paranoidDir, realPath string) (lastNode string, paranoidPath string) {
	split := strings.Split(realPath, "/")
	paranoidPath = path.Join(paranoidDir, "names")
	lastNode = ""
	for i := range split {
		current := split[i] + "-file"
		paranoidPath = path.Join(paranoidPath, current)
		if i == len(split)-1 {
			lastNode = current
		}
	}
	return lastNode, paranoidPath
}

func generateNewInode() (inodeBytes []byte, inodeString string) {
	inodeBytes, err := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	checkErr("util, generateNewInode", err)
	inodeString = strings.TrimSpace(string(inodeBytes))
	return []byte(inodeString), inodeString
}

func getFileInode(filePath string) (inodeBytes []byte, inodeString string, errorCode int) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", returncodes.ENOENT
		} else if os.IsPermission(err) {
			return nil, "", returncodes.EACCES
		}
	}
	return bytes, string(bytes), returncodes.OK
}

func deleteFile(filePath string) (returncode int) {
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return returncodes.ENOENT
		} else if os.IsPermission(err) {
			return returncodes.EACCES
		}
	}
	return returncodes.OK
}

func isDirectory(path string) bool {
	f, err := os.Open(path)
	checkErr("util, isDirectory", err)
	defer f.Close()
	fi, err := f.Stat()
	checkErr("util, isDirectory", err)
	return fi.Mode().IsDir()
}
