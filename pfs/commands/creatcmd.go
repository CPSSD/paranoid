package commands

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type inode struct {
	Count int    `json:"count"`
	Inode string `json:"inode"`
}

func CreatCommand(args []string) {
	if len(args) < 2 {
		log.Fatal("Not enough arguments!")
	}
	directory := args[0]
	if _, err := os.Stat(path.Join(directory, "/names/", args[1])); os.IsNotExist(err) == false {
		log.Fatal("creat : file already exits")
	}
	uuidbytes, err := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	uuid := strings.TrimSpace(string(uuidbytes))
	checkErr("creat", err)
	err = ioutil.WriteFile(path.Join(directory, "/names/", args[1]), []byte(uuid), 0777)
	checkErr("creat", err)
	nodeData := &inode{
		Inode: uuid,
		Count: 1}
	jsonData, err := json.Marshal(nodeData)
	checkErr("creat", err)
	err = ioutil.WriteFile(path.Join(directory, "/inodes/", uuid), jsonData, 0777)
	checkErr("creat", err)
	_, err = os.Create(path.Join(directory, "/contents/", uuid))
	checkErr("creat", err)
}
