package commands

import (
	"github.com/cpssd/paranoid/pfsm/icclient"
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
	paranoidPath = ""
	lastNode = ""
	for i := range split {
		paranoidPath += split[i] + "-file"
		if i != len(split)-1 {
			paranoidPath += "/"
		} else {
			lastNode = split[i] + "-file"
		}
	}
	return lastNode, path.Join(paranoidDir, "names", paranoidPath)
}

func generateNewInode() ([]byte, string) {
	uuidbytes, err := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	checkErr("util, generateNewInode", err)
	uuidString := strings.TrimSpace(string(uuidbytes))
	return []byte(uuidString), uuidString
}
