package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"os"
	"path"
	"syscall"
	"time"
)

type statInfo struct {
	Length int64
	Ctime  time.Time
	Mtime  time.Time
	Atime  time.Time
	Mode   os.FileMode
}

//StatCommand returns information about a file
func StatCommand(paranoidDirectory, filePath string) (returnCode int, returnError error, info statInfo) {
	Log.Verbose("stat : given paranoidDirectory", paranoidDirectory)

	err := getFileSystemLock(paranoidDirectory, sharedLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, statInfo{}
	}

	defer func() {
		err := unLockFileSystem(paranoidDirectory)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
			info = statInfo{}
		}
	}()

	namepath := getParanoidPath(paranoidDirectory, filePath)
	namePathType, err := getFileType(paranoidDirectory, namepath)
	if err != nil {
		return returncodes.EUNEXPECTED, err, statInfo{}
	}

	if namePathType == typeENOENT {
		return returncodes.ENOENT, errors.New(filePath + " does not exist"), statInfo{}
	}

	inodeBytes, code, err := getFileInode(namepath)
	if code != returncodes.OK {
		return code, err, statInfo{}
	}

	inodeName := string(inodeBytes)
	contentsFile := path.Join(paranoidDirectory, "contents", inodeName)

	fi, err := os.Lstat(contentsFile)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error Lstating file:", err), statInfo{}
	}

	stat := fi.Sys().(*syscall.Stat_t)
	atime := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
	ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))

	var mode os.FileMode
	switch namePathType {
	case typeDir:
		mode = os.FileMode(syscall.S_IFDIR | fi.Mode().Perm())
	default:
		mode = os.FileMode(stat.Mode)
	}

	statData := &statInfo{
		Length: fi.Size(),
		Mtime:  fi.ModTime(),
		Ctime:  ctime,
		Atime:  atime,
		Mode:   mode}

	Log.Verbose("stat : returning", statData)
	return returncodes.OK, nil, *statData
}
