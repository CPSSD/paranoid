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
	verboseLog("readdir : given directory = " + directory)
	files, err := ioutil.ReadDir(path.Join(directory, "names"))
	checkErr("readdir", err)
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
	for i := 0; i < len(files); i++ {
		fmt.Println(files[i].Name())
	}
}
