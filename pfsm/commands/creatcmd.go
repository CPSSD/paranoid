package commands

import (
	"encoding/json"
	"github.com/cpssd/paranoid/pfsm/network"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

type inode struct {
	Count int    `json:"count"`
	Inode string `json:"inode"`
}

//CreatCommand creates a new file with the name args[1] in the pfs directory args[0]
func CreatCommand(args []string) {
	verboseLog("creat command called")
	if len(args) < 3 {
		log.Fatalln("Not enough arguments!")
	}
	directory := args[0]
	verboseLog("creat : directory = " + directory)
	if _, err := os.Stat(path.Join(directory, "names", args[1])); !os.IsNotExist(err) {
		log.Fatalln("creat : file already exits")
	}
	verboseLog("creat : creating file " + args[1])
	uuidbytes, err := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	checkErr("creat", err)
	uuid := strings.TrimSpace(string(uuidbytes))
	verboseLog("creat : uuid = " + uuid)
	perms, err := strconv.Atoi(args[2])
	checkErr("creat", err)
	err = ioutil.WriteFile(path.Join(directory, "names", args[1]), []byte(uuid), os.FileMode(perms))
	checkErr("creat", err)
	nodeData := &inode{
		Inode: uuid,
		Count: 1}
	jsonData, err := json.Marshal(nodeData)
	checkErr("creat", err)
	err = ioutil.WriteFile(path.Join(directory, "inodes", uuid), jsonData, 0777)
	checkErr("creat", err)
	_, err = os.Create(path.Join(directory, "contents", uuid))
	checkErr("creat", err)
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
	if !Flags.Network {
		network.Creat(directory, args[1])
	}
}
