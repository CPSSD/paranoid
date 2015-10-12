package commands

import (
	"io/ioutil"
	"log"
	"path"
)

func MountCommand(args []string) {
	if len(args) < 3 {
		log.Fatal("Not enough arguments!")
	}
	directory := args[0]
	err := ioutil.WriteFile(path.Join(directory, "meta/ip"), []byte(args[1]), 0777)
	checkErr(err)
	err = ioutil.WriteFile(path.Join(directory, "meta/port"), []byte(args[2]), 0777)
	checkErr(err)
}
