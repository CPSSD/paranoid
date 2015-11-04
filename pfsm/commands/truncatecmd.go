package commands

import (
	"github.com/cpssd/paranoid/pfsm/network"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

//TruncateCommand reduces the file given as args[1] in the paranoid-direcory args[0] to the size given in args[2]
func TruncateCommand(args []string) {
	verboseLog("truncate command given")
	if len(args) < 3 {
		log.Fatalln("Not enough arguments!")
	}
	directory := args[0]
	verboseLog("truncate : given directory = " + directory)
	if !checkFileExists(path.Join(directory, "names", args[1])) {
		io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.ENOENT))
		return
	}
	fileNameBytes, err := ioutil.ReadFile(path.Join(directory, "names", args[1]))
	checkErr("truncate", err)
	fileName := string(fileNameBytes)
	verboseLog("truncate : truncating " + fileName)
	newsize, err := strconv.Atoi(args[2])
	checkErr("truncate", err)
	contentsFile, err := os.OpenFile(path.Join(directory, "contents", fileName), os.O_WRONLY, 0777)
	checkErr("truncate", err)
	err = contentsFile.Truncate(int64(newsize))
	checkErr("truncate", err)
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
	if !Flags.Network {
		network.Truncate(directory, args[1], newsize)
	}
}
