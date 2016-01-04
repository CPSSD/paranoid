package commands

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"io/ioutil"
	"path"
	"strings"
)

//ReadDirCommand returns a list of all the files in the given directory
func ReadDirCommand(directory, dirName string) (returnCode int, returnError error, fileNames []string) {
	Log.Info("readdir command called")
	Log.Verbose("readdir : given directory = " + directory)

	err := getFileSystemLock(directory, sharedLock)
	if err != nil {
		return returncodes.EUNEXPECTED, err, nil
	}

	defer func() {
		err := unLockFileSystem(directory)
		if err != nil {
			returnCode = returncodes.EUNEXPECTED
			returnError = err
			fileNames = nil
		}
	}()

	dirpath := ""

	if dirName == "" {
		dirpath = path.Join(directory, "names")
	} else {
		dirpath = getParanoidPath(directory, dirName)
		pathFileType, err := getFileType(directory, dirpath)
		if err != nil {
			return returncodes.EUNEXPECTED, err, nil
		}

		if pathFileType == typeENOENT {
			return returncodes.ENOENT, errors.New(dirName + " does not exist"), nil
		}

		if pathFileType == typeFile {
			return returncodes.ENOTDIR, errors.New(dirName + " is of type file"), nil
		}

		if pathFileType == typeSymlink {
			return returncodes.ENOTDIR, errors.New(dirName + " is of type symlink"), nil
		}
	}

	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return returncodes.EUNEXPECTED, fmt.Errorf("error reading directory "+dirName+":", err), nil
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
