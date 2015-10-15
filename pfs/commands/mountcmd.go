package commands

import (
	"io/ioutil"
	"log"
	"path"
)

func MountCommand(args []string) {
	verboseLog("mount command called")
	if len(args) < 3 {
		log.Fatal("Not enough arguments!")
	}
	directory := args[0]
	verboseLog("mount : given directory = " + directory)
	err := ioutil.WriteFile(path.Join(directory, "meta", "ip"), []byte(args[1]), 0777)
	checkErr("mount", err)
	err = ioutil.WriteFile(path.Join(directory, "meta", "port"), []byte(args[2]), 0777)
	checkErr("mount", err)
}
