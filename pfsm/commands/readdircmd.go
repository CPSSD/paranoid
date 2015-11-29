package commands

import (
	"fmt"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

//ReadDirCommand takes a pfs directory as args[0] and prints a list of the names of the files in that directory 1 per line.
func ReadDirCommand(args []string) {
	verboseLog("readdir command called")
	if len(args) < 1 {
		log.Fatalln("Not enough arguments!")
	}

	directory := args[0]
	getFileSystemLock(directory, sharedLock)
	defer unLockFileSystem(directory)
	verboseLog("readdir : given directory = " + directory)

	dirpath := ""
	dirInfoName := ""

	if args[1] == directory {
		dirpath = path.Join(directory, "names")
	} else {
		dirname, dirp := getParanoidPath(directory, args[1])
		dirpath = dirp
		dirInfoName = dirname + "-info"
	}

	files, err := ioutil.ReadDir(dirpath)
	checkErr("readdir", err)
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
	for i := 0; i < len(files); i++ {
		file := files[i].Name()
		if dirInfoName != "" {
			if file != dirInfoName {
				fmt.Println(files[i].Name())
			}
		} else {
			fmt.Println(files[i].Name())
		}
	}
}
