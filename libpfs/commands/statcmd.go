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

// StatCommand returns information about a file as StatInfo object
func StatCommand(paranoidDirectory, filePath string) (returnCode returncodes.Code, returnError error, info statInfo) {
	Log.Info("stat command called")
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
	contentsFilePath := path.Join(paranoidDirectory, "contents", inodeName)

	contentsFile, err := os.Open(contentsFilePath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error opening contents file: %s", err), statInfo{}
	}

	fi, err := os.Lstat(contentsFilePath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error Lstating file: %s", err), statInfo{}
	}

	stat := fi.Sys().(*syscall.Stat_t)
	atime := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
	ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
	mode, err := getFileMode(paranoidDirectory, inodeName)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error getting filemode: %s", err), statInfo{}
	}

	fileLength, err := getFileLength(contentsFile)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error getting file length: %s", err), statInfo{}
	}

	statData := &statInfo{
		Length: fileLength,
		Mtime:  fi.ModTime(),
		Ctime:  ctime,
		Atime:  atime,
		Mode:   mode}

	Log.Verbose("stat : returning", statData)
	return returncodes.OK, nil, *statData
}
