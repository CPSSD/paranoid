package mountcmd

import (
	"io/ioutil"
	"log"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal("mount error occured:", err)
	}
}

func MountCommand(args []string) {
	if len(args) < 3 {
		log.Fatal("Not enough arguments!")
	}
	directory := args[0]
	if directory[len(directory)-1] != '/' {
		directory += "/"
	}
	err := ioutil.WriteFile(directory+"meta/ip", []byte(args[1]), 0777)
	checkErr(err)
	err = ioutil.WriteFile(directory+"meta/port", []byte(args[2]), 0777)
	checkErr(err)
}
