package commands

import (
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"log"
	"os"
	"path"
)

// RenameCommand is called when renaming a file
func RenameCommand(args []string) {
	verboseLog("rename command called")
	if len(args) < 3 {
		log.Fatalln("Not enough arguments!")
	}

	directory := args[0]
	oldFilePath := path.Join(directory, "names", args[1])
	newFilePath := path.Join(directory, "names", args[2])

	if !checkFileExists(oldFilePath) {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}

	err := os.Rename(oldFilePath, newFilePath)
	checkErr("rename", err)

	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
