package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"io/ioutil"
	"path"
	"strings"
)

//ReadDirCommand returns a list of all the files in the given paranoidDirectory
func ReadDirCommand(paranoidDirectory, dirPath string) (returnCode int, returnError error, fileNames []string) {
	Log.Info("readdir command called")
	Log.Verbose("readdir : given paranoidDirectory = " + paranoidDirectory)

	err := getFileSystemLock(paranoidDirectory, sharedLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, nil
	}

	defer func() {
		err := unLockFileSystem(paranoidDirectory)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
			fileNames = nil
		}
	}()

	dirParanoidPath := ""

	if dirPath == "" {
		dirParanoidPath = path.Join(paranoidDirectory, "names")
	} else {
		dirParanoidPath = getParanoidPath(paranoidDirectory, dirPath)
		pathFileType, err := getFileType(paranoidDirectory, dirParanoidPath)
		if err != nil {
			return returncodes.EUNEXPECTED, err, nil
		}

		if pathFileType == typeENOENT {
			return returncodes.ENOENT, errors.New(dirPath + " does not exist"), nil
		}

		if pathFileType == typeFile {
			return returncodes.ENOTDIR, errors.New(dirPath + " is of type file"), nil
		}

		if pathFileType == typeSymlink {
			return returncodes.ENOTDIR, errors.New(dirPath + " is of type symlink"), nil
		}
	}

	files, err := ioutil.ReadDir(dirParanoidPath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error reading paranoidDirectory "+dirPath+":", err), nil
	}

	var names []string
	for i := 0; i < len(files); i++ {
		file := files[i].Name()
		if file != "info" {
			names = append(names, file[:strings.LastIndex(file, "-")])
		}
	}
	return returncodes.OK, nil, names
}
