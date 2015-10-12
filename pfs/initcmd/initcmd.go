package initcmd

import (
	"io/ioutil"
	"log"
	"os"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal("init error occured:", err)
	}
}

func makeDir(parentDir, newDir string) string {
	newDirPath := parentDir
	if parentDir[len(parentDir)-1] != '/' {
		newDirPath += "/"
	}
	newDirPath += newDir
	err := os.Mkdir(newDirPath, 0777)
	checkErr(err)
	return newDirPath
}

func checkEmpty(directory string) {
	files, err := ioutil.ReadDir(directory)
	checkErr(err)
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
	makeDir(directory, "names/")
	makeDir(directory, "inodes/")
	metaDir := makeDir(directory, "meta/")
	makeDir(directory, "contents/")
	uuid, err := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	checkErr(err)
	err = ioutil.WriteFile(metaDir+"uuid", uuid, 0777)
	checkErr(err)
}
