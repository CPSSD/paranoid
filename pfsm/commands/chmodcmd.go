package commands

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

//ChmodCommand is used to change the permissions of a file.
func ChmodCommand(args []string) {
	verboseLog("chmod command given")
	if len(args) < 3 {
		log.Fatalln("Not enough arguments!")
	}
	directory := args[0]
	verboseLog("chmod : given directory = " + directory)
	if !checkFileExists(path.Join(directory, "names", args[1])) {
		io.WriteString(os.Stdout, getReturnCode(ENOENT))
		return
	}
	fileNameBytes, err := ioutil.ReadFile(path.Join(directory, "names", args[1]))
	checkErr("chmod", err)
	fileName := string(fileNameBytes)
	verboseLog("chmod : changing permissions of " + fileName + " to " + args[2])
	perms, err := strconv.Atoi(args[2])
	checkErr("chmod", err)
	contentsFile, err := os.OpenFile(path.Join(directory, "contents", fileName), os.O_WRONLY, 0777)
	checkErr("chmod", err)
	err = contentsFile.Chmod(os.FileMode(perms))
	checkErr("chmod", err)
	io.WriteString(os.Stdout, getReturnCode(OK))
}
