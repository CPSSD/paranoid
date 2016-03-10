package commands

import (
	"fmt"
	"github.com/cpssd/paranoid/libpfs"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"os"
	"path"
	"syscall"
	"time"
)

// StatInfo contains the stat information of the file
type StatInfo struct {
	Length int64
	Ctime  time.Time
	Mtime  time.Time
	Atime  time.Time
	Mode   os.FileMode
}

//StatCommand returns information about a file
func StatCommand(paranoidDirectory, filePath string) (returnCode int, returnError error, info statInfo) {
	Log.Info("stat command called")
	Log.Verbose("stat : given paranoidDirectory", paranoidDirectory)

	err := getFileSystemLock(paranoidDirectory, sharedLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, StatInfo{}
	}

	defer func() {
		err := unLockFileSystem(paranoidDirectory)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
			info = StatInfo{}
		}
	}()

	namepath := getParanoidPath(paranoidDirectory, filePath)
	namePathType, err := getFileType(paranoidDirectory, namepath)
	if err != nil {
		return returncodes.EUNEXPECTED, err, StatInfo{}
	}

	if namePathType == typeENOENT {
		return returncodes.ENOENT, fmt.Errorf("%s does not exist", filePath), StatInfo{}
	}

	inodeBytes, code, err := getFileInode(namepath)
	if code != returncodes.OK {
		return code, err, StatInfo{}
	}

	inodeName := string(inodeBytes)
	contentsFile := path.Join(paranoidDirectory, "contents", inodeName)

	fi, err := os.Lstat(contentsFile)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error Lstating file: %s", err), StatInfo{}
	}

	stat := fi.Sys().(*syscall.Stat_t)
	atime := time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
	ctime := time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
	mode, err := getFileMode(paranoidDirectory, inodeName)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error getting filemode:", err), statInfo{}
	}

	// Get the size of the last last block
	file, err := os.OpenFile(contentsFile, 0700, mode)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error opening file: %s", err), StatInfo{}
	}
	defer file.Close()
	finalBlockSize, err := libpfs.LastBlockSize(file)

	// Return the file with correction for the size
	statData := StatInfo{
		Length: fi.Size() - int64(libpfs.CipherBlock.BlockSize()-finalBlockSize),
		Mtime:  fi.ModTime(),
		Ctime:  ctime,
		Atime:  atime,
		Mode:   mode,
	}

	Log.Verbose("stat: returning", statData)
	return returncodes.OK, nil, statData
}
