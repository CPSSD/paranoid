package commands

import (
	"fmt"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

//ReadDirCommand takes a pfs directory as args[0] and prints a list of the names of the files in that directory 1 per line.
func ReadDirCommand(args []string) {
	verboseLog("readdir command called")
	if len(args) < 2 {
		log.Fatalln("Not enough arguments!")
	}

	directory := args[0]
	getFileSystemLock(directory, sharedLock)
	defer unLockFileSystem(directory)
	verboseLog("readdir : given directory = " + directory)

	dirpath := ""
	dirInfoName := ""

	if args[1] == "" {
		dirpath = path.Join(directory, "names")
	} else {
		dirp := getParanoidPath(directory, args[1])
		dirpath = dirp
		dirInfoName = path.Base(dirp) + "-info"
	}

	files, err := ioutil.ReadDir(dirpath)
	checkErr("readdir", err)
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
	for i := 0; i < len(files); i++ {
		file := files[i].Name()
		realFileName := file[:strings.LastIndex(file, "-")]
		if dirInfoName != "" {
			if file != dirInfoName {
				fmt.Println(realFileName)
			}
		} else {
			fmt.Println(realFileName)
		}
	}
}
