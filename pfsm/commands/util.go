package commands

import (
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/pfsm/icclient"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"syscall"
)

var Log *logger.ParanoidLogger

type inode struct {
	Count int    `json:"count"`
	Inode string `json:"inode"`
	Link  string `json:"link,omitempty"`
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

func getParanoidPath(paranoidDir, realPath string) (paranoidPath string) {
	split := strings.Split(realPath, "/")
	paranoidPath = path.Join(paranoidDir, "names")
	for i := range split {
		paranoidPath = path.Join(paranoidPath, (split[i] + "-file"))
	}
	return paranoidPath
}

func generateNewInode() (inodeBytes []byte) {
	inodeBytes, err := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	if err != nil {
		Log.Fatal("error generating new Inode:", err)
	}
	return []byte(strings.TrimSpace(string(inodeBytes)))
}

func getFileInode(filePath string) (inodeBytes []byte, errorCode int) {
	f, err := os.Lstat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, returncodes.ENOENT
		}
		Log.Fatal("util, getFileInode", " error occured: ", err)
	}
	if f.Mode().IsDir() {
		filePath = path.Join(filePath, "info")
	}
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, returncodes.ENOENT
		} else if os.IsPermission(err) {
			return nil, returncodes.EACCES
		}
		log.Fatalln("util, getFileInode", " error occured: ", err)
	}
	return bytes, returncodes.OK
}

func deleteFile(filePath string) (returncode int) {
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return returncodes.ENOENT
		} else if os.IsPermission(err) {
			return returncodes.EACCES
		}
		log.Fatalln("util, deleteFile", " error occured: ", err)
	}
	return returncodes.OK
}

const (
	typeFile = iota
	typeDir
	typeSymlink
	typeENOENT
)

func getFileType(directory, filePath string) int {
	f, err := os.Lstat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return typeENOENT
		}
		Log.Fatal("util, getFileType", " error occured: ", err)
	}
	if f.Mode().IsDir() {
		return typeDir
	}

	inode, errorcode := getFileInode(filePath)
	if errorcode != returncodes.OK {
		if errorcode == returncodes.ENOENT {
			return typeENOENT
		}
		Log.Fatal("util, getFileType symlink check error occured code: ", errorcode)
	}
	f, err = os.Lstat(path.Join(directory, "contents", string(inode)))
	if err != nil {
		Log.Fatal("util, getFileType symlink check error occured: ", err)
	}

	if f.Mode()&os.ModeSymlink > 0 {
		log.Println("Is symlink")
		return typeSymlink
	}

	return typeFile
}
