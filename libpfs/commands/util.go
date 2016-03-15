package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/logger"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"syscall"
)

var Log *logger.ParanoidLogger

type inode struct {
	Count int         `json:"count"`
	Inode string      `json:"inode"`
	Mode  os.FileMode `json:"mode"`
	Link  string      `json:"link,omitempty"`
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

func getFileMode(paranoidDir, inodeName string) (os.FileMode, error) {
	inodePath := path.Join(paranoidDir, "inodes", inodeName)
	inodeContents, err := ioutil.ReadFile(inodePath)
	if err != nil {
		return os.FileMode(0), fmt.Errorf("error reading inode: %s", err)
	}

	nodeData := &inode{}
	err = json.Unmarshal(inodeContents, &nodeData)
	if err != nil {
		return os.FileMode(0), fmt.Errorf("error unmarshaling inode data: %s", err)
	}
	return nodeData.Mode, nil
}

func canAccessFile(paranoidDir, inodeName string, mode uint32) (int, error) {
	fileMode, err := getFileMode(paranoidDir, inodeName)
	if err != nil {
		return returncodes.EUNEXPECTED, err
	}

	userPerms := (uint32(fileMode) >> 6) & 7
	if userPerms&mode != mode {
		return returncodes.EACCES, errors.New("invalid permissions")
	}
	return returncodes.OK, nil
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

func getParanoidPath(paranoidDir, realPath string) (paranoidPath string) {
	split := strings.Split(realPath, "/")
	paranoidPath = path.Join(paranoidDir, "names")
	for i := range split {
		paranoidPath = path.Join(paranoidPath, (split[i] + "-file"))
	}
	return paranoidPath
}

func generateNewInode() (inodeBytes []byte, err error) {
	inodeBytes, err = ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	if err != nil {
		return nil, fmt.Errorf("error generating new inode:", err)
	}
	return []byte(strings.TrimSpace(string(inodeBytes))), nil
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

func deleteFile(filePath string) (returncode int, returnerror error) {
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return returncodes.ENOENT, errors.New("error deleting file: " + filePath + " does not exist")
		} else if os.IsPermission(err) {
			return returncodes.EACCES, errors.New("error deleting file: could not access " + filePath)
		}
		return returncodes.EUNEXPECTED, fmt.Errorf("unexpected error deleting file: ", err)
	}
	return returncodes.OK, nil
}

const (
	typeFile = iota
	typeDir
	typeSymlink
	typeENOENT
)

func getFileType(directory, filePath string) (int, error) {
	f, err := os.Lstat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return typeENOENT, nil
		}
		return 0, fmt.Errorf("error getting file type of "+filePath+", error stating file:", err)
	}

	if f.Mode().IsDir() {
		return typeDir, nil
	}

	inode, code, err := getFileInode(filePath)
	if err != nil {
		return 0, fmt.Errorf("error getting file type of "+filePath+", error getting inode:", err)
	}

	if code != returncodes.OK {
		if code == returncodes.ENOENT {
			return typeENOENT, nil
		}
		return 0, fmt.Errorf("error getting file type of "+filePath+", unexpected result from getFileInode:", code)
	}

	f, err = os.Lstat(path.Join(directory, "contents", string(inode)))
	if err != nil {
		return 0, fmt.Errorf("error getting file type of "+filePath+", symlink check error occured:", err)
	}

	if f.Mode()&os.ModeSymlink > 0 {
		return typeSymlink, nil
	}

	return typeFile, nil
}
