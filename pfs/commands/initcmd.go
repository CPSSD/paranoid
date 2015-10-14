package commands

import (
	"io/ioutil"
	"log"
	"os"
	"path"
)

func checkErr(cmd string, err error) {
	if err != nil {
		log.Fatal(cmd, " error occured: ", err)
	}
}

func makeDir(parentDir, newDir string) string {
	newDirPath := path.Join(parentDir, newDir)
	err := os.Mkdir(newDirPath, 0777)
	checkErr("init", err)
	return newDirPath
}

func checkEmpty(directory string) {
	files, err := ioutil.ReadDir(directory)
	checkErr("init", err)
	if len(files) > 0 {
		log.Fatal("init : directory must be empty")
	}
}

func InitCommand(args []string) {
	if len(args) < 1 {
		log.Fatal("Not enough arguments!")
	}
	directory := args[0]
	checkEmpty(directory)
	makeDir(directory, "names")
	makeDir(directory, "inodes")
	metaDir := makeDir(directory, "meta")
	makeDir(directory, "contents")
	uuid, err := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	checkErr("init", err)
	err = ioutil.WriteFile(path.Join(metaDir, "uuid"), uuid, 0777)
	checkErr("init", err)
}
