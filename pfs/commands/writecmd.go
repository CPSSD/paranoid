package commands

import (
	"github.com/cpssd/paranoid/pfs/network"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

//WriteCommand writes data from Stdin to the file given in args[1] in the pfs directory args[0]
//Can also be given an offset and length as args[2] and args[3] otherwise it writes from the start of the file
func WriteCommand(args []string) {
	verboseLog("write command given")
	if len(args) < 2 {
		log.Fatalln("Not enough arguments!")
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
		network.Write(directory, fileName, -1, -1, fileData)
	} else {
		contentsFile, err := os.OpenFile(path.Join(directory, "contents", fileName), os.O_WRONLY, 0777)
		checkErr("write", err)
		offset, err := strconv.Atoi(args[2])
		checkErr("write", err)
		if len(args) == 3 {
			err = contentsFile.Truncate(int64(offset))
			checkErr("write", err)
		} else {
			length, err := strconv.Atoi(args[3])
			checkErr("write", err)
			if len(fileData) > length {
				fileData = fileData[:length]
			} else if len(fileData) < length {
				emptyBytes := make([]byte, length-len(fileData))
				fileData = append(fileData, emptyBytes...)
			}
		}
		_, err = contentsFile.WriteAt(fileData, int64(offset))
		checkErr("write", err)
		if len(args) == 3 {
			network.Write(directory, fileName, offset, -1, fileData)
		} else {
			network.Write(directory, fileName, offset, length, fileData)
		}
	}
}
