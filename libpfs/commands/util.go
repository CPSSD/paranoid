package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/encryption"
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

func canAccessFile(paranoidDir, inodeName string, mode uint32) (returncodes.Code, error) {
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

func getFileLength(file *os.File) (int64, error) {
	fi, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("unable to get file length: %s", err)
	}
	if !encryption.Encrypted {
		return fi.Size(), nil
	}
	if fi.Size() == 0 {
		return 0, nil
	}
	lastBlockSize, err := encryption.LastBlockSize(file)
	if err != nil {
		return 0, fmt.Errorf("unable to get file length: %s", err)
	}
	fileLength := fi.Size() - 1
	if fileLength > 0 {
		fileLength = fileLength - int64(encryption.GetCipherSize()-lastBlockSize)
	}
	return fileLength, nil
}

//Types of locks
const (
	SharedLock = iota
	ExclusiveLock
)

func GetFileSystemLock(paranoidDir string, lockType int) error {
	lockPath := path.Join(paranoidDir, "meta", "lock")
	file, err := os.Open(lockPath)
	if err != nil {
		return fmt.Errorf("could not get meta/lock file descriptor: %s", err)
	}
	defer file.Close()
	if lockType == SharedLock {
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_SH)
		if err != nil {
			return fmt.Errorf("error getting shared lock on meta/lock: %s", err)
		}
	} else if lockType == ExclusiveLock {
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
		if err != nil {
			return fmt.Errorf("error getting exclusive lock on meta/lock: %s", err)
		}
	}
	return nil
}

func getFileLock(paranoidDir, fileName string, lockType int) error {
	lockPath := path.Join(paranoidDir, "contents", fileName)
	file, err := os.Open(lockPath)
	if err != nil {
		return fmt.Errorf("could not get file descriptor for lock file %s: %s", fileName, err)
	}
	defer file.Close()
	if lockType == SharedLock {
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_SH)
		if err != nil {
			return fmt.Errorf("could not get shared lock on lock file %s: %s", fileName, err)
		}
	} else if lockType == ExclusiveLock {
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
		if err != nil {
			return fmt.Errorf("could not get exclusive lock on lock file %s: %s", fileName, err)
		}
	}
	return nil
}

func UnLockFileSystem(paranoidDir string) error {
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
		return fmt.Errorf("could not get file descriptor for unlock file %s: %s", fileName, err)
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
		return nil, fmt.Errorf("error generating new inode: %s", err)
	}
	return []byte(strings.TrimSpace(string(inodeBytes))), nil
}

func getFileInode(filePath string) (inodeBytes []byte, errorCode returncodes.Code, err error) {
	f, err := os.Lstat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, returncodes.ENOENT, errors.New("error getting inode, " + filePath + " does not exist")
		}
		return nil, returncodes.EUNEXPECTED, fmt.Errorf("unexpected error getting inode of file %s: %s", filePath, err)
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
		return nil, returncodes.EUNEXPECTED, fmt.Errorf("unexpected error getting inode of file %s: %s", filePath, err)
	}
	return bytes, returncodes.OK, nil
}

func deleteFile(filePath string) (returncode returncodes.Code, returnerror error) {
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return returncodes.ENOENT, errors.New("error deleting file: " + filePath + " does not exist")
		} else if os.IsPermission(err) {
			return returncodes.EACCES, errors.New("error deleting file: could not access " + filePath)
		}
		return returncodes.EUNEXPECTED, fmt.Errorf("unexpected error deleting file: %s", err)
	}
	return returncodes.OK, nil
}

const (
	typeFile = iota
	typeDir
	typeSymlink
	typeENOENT
)

func getFileType(directory, filePath string) (returncodes.Code, error) {
	f, err := os.Lstat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return typeENOENT, nil
		}
		return 0, fmt.Errorf("error getting file type of %s, error stating file: %s", filePath, err)
	}

	if f.Mode().IsDir() {
		return typeDir, nil
	}

	inode, code, err := getFileInode(filePath)
	if err != nil {
		return 0, fmt.Errorf("error getting file type of %s, error getting inode: %s", filePath, err)
	}

	if code != returncodes.OK {
		if code == returncodes.ENOENT {
			return typeENOENT, nil
		}
		return 0, fmt.Errorf("error getting file type of %s, unexpected result from getFileInode: %s", filePath, code)
	}

	f, err = os.Lstat(path.Join(directory, "contents", string(inode)))
	if err != nil {
		return 0, fmt.Errorf("error getting file type of %s, symlink check error occured: %s", filePath, err)
	}

	if f.Mode()&os.ModeSymlink > 0 {
		return typeSymlink, nil
	}

	return typeFile, nil
}
