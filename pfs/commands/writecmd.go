package commands

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

//WriteCommand writes data from Stdin to the file given in args[1] in the pfs directory args[0]
//Can also be given an offset and length as args[2] and args[3]
func WriteCommand(args []string) {
	verboseLog("write command given")
	if len(args) < 2 {
		log.Fatal("Not enough arguments!")
	}
	directory := args[0]
	verboseLog("write : given directory = " + directory)
	fileNameBytes, err := ioutil.ReadFile(path.Join(directory, "names", args[1]))
	checkErr("write", err)
	fileName := string(fileNameBytes)
	verboseLog("write : wrting to " + fileName)
	fileData, err := ioutil.ReadAll(os.Stdin)
	checkErr("write", err)
	if len(args) == 2 {
		err = ioutil.WriteFile(path.Join(directory, "contents", fileName), fileData, 0777)
		checkErr("write", err)
	} else {
		contentsFile, err := os.OpenFile(path.Join(directory, "contents", fileName), os.O_WRONLY, 0777)
		checkErr("write", err)
		offset, err := strconv.Atoi(args[2])
		checkErr("write", err)
		_, err = contentsFile.WriteAt(fileData, int64(offset))
		checkErr("write", err)
	}
}
