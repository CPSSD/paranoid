package commands

import (
	"errors"
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

func getFileSystemLock(paranoidDir string, lockType int) error {
	lockPath := path.Join(paranoidDir, "meta", "lock")
	file, err := os.Open(lockPath)
	if err != nil {
		return fmt.Errorf("could not get meta/lock file descriptor:", err)
	}
	defer file.Close()
	if lockType == sharedLock {
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_SH)
		if err != nil {
			return fmt.Errorf("error getting shared lock on meta/lock:", err)
		}
	} else if lockType == exclusiveLock {
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
		if err != nil {
			return fmt.Errorf("error getting exclusive lock on meta/lock:", err)
		}
	}
	return nil
}

func getFileLock(paranoidDir, fileName string, lockType int) error {
	lockPath := path.Join(paranoidDir, "contents", fileName)
	file, err := os.Open(lockPath)
	if err != nil {
		return fmt.Errorf("could not get file descriptor for lock file "+fileName+" :", err)
	}
	defer file.Close()
	if lockType == sharedLock {
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_SH)
		if err != nil {
			return fmt.Errorf("could not get shared lock on lock file "+fileName+" :", err)
		}
	} else if lockType == exclusiveLock {
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
		if err != nil {
			return fmt.Errorf("could not get exclusive lock on lock file "+fileName+" :", err)
		}
	}
	return nil
}

func unLockFileSystem(paranoidDir string) error {
	lockPath := path.Join(paranoidDir, "meta", "lock")
	file, err := os.Open(lockPath)
	if err != nil {
		return errors.New("could not get meta/lock file descriptor for unlock")
	}
	defer file.Close()
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
	if err != nil {
		return errors.New("could not unlock meta/lock")
	}
	return nil
}

func unLockFile(paranoidDir, fileName string) error {
	lockPath := path.Join(paranoidDir, "contents", fileName)
	file, err := os.Open(lockPath)
	if err != nil {
		return fmt.Errorf("could not get file descriptor for unlock file "+fileName+" :", err)
	}
	defer file.Close()
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
	if err != nil {
		return errors.New("could not unlock " + fileName)
	}
	return nil
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

func generateNewInode() (inodeBytes []byte, err error) {
	inodeBytes, err := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	if err != nil {
		return nil, fmt.Errorf("error generating new inode:", err)
	}
	return []byte(strings.TrimSpace(string(inodeBytes)))
}

func getFileInode(filePath string) (inodeBytes []byte, errorCode int, err error) {
	f, err := os.Lstat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, returncodes.ENOENT, errors.New("error getting inode, " + filePath + " does not exist")
		}
		return nil, returncodes.EUNEXPECTED, fmt.Errorf("unexpected error getting inode of file "+filePath+" :", err)
	}
	if f.Mode().IsDir() {
		filePath = path.Join(filePath, "info")
	}
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, returncodes.ENOENT, errors.New("error getting inode, " + filePath + "does not exist")
		} else if os.IsPermission(err) {
			return nil, returncodes.EACCES, errors.New("error getting inode, could not access " + filePath)
		}
		return nil, returncodes.EUNEXPECTED, fmt.Errorf("unexpected error getting inode of file "+filePath+" :", err)
	}
	return bytes, returncodes.OK, nil
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
