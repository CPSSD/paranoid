package commands

import (
	"fmt"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

//ReadDirCommand takes a pfs directory as args[0] and prints a list of the names of the files in that directory 1 per line.
func ReadDirCommand(args []string) {
	Log.Info("readdir command called")
	if len(args) < 2 {
		Log.Fatal("Not enough arguments!")
	}

	directory := args[0]
	getFileSystemLock(directory, sharedLock)
	defer unLockFileSystem(directory)
	Log.Verbose("readdir : given directory = " + directory)

	dirpath := ""

	if args[1] == "" {
		dirpath = path.Join(directory, "names")
	} else {
		dirpath = getParanoidPath(directory, args[1])
		pathFileType := getFileType(directory, dirpath)
		if pathFileType == typeENOENT {
			io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
			return
		}
		if pathFileType == typeFile {
			io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOTDIR))
			return
		}
		if pathFileType == typeSymlink {
			io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.EIO))
			return
		}
	}

	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		Log.Fatal("error reading directory:", err)
	}

	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
	for i := 0; i < len(files); i++ {
		file := files[i].Name()
		if file != "info" {
			fmt.Println(file[:strings.LastIndex(file, "-")])
		}
	}
}
