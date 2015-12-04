package commands

import (
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

//makeDir creates a new directory with permissions 0777 with the name newDir in parentDir.
func makeDir(parentDir, newDir string) string {
	newDirPath := path.Join(parentDir, newDir)
	err := os.Mkdir(newDirPath, 0700)
	checkErr("init", err)
	return newDirPath
}

//checkEmpty checks if a given directory has any children.
func checkEmpty(directory string) {
	files, err := ioutil.ReadDir(directory)
	checkErr("init", err)
	if len(files) > 0 {
		log.Fatalln("init : directory must be empty")
	}
}

//InitCommand creates the pvd directory sturucture in args[0]
//It also gets a UUID and stores it in the meta directory.
func InitCommand(args []string) {
	Log.Verbose("init command called")
	if len(args) < 1 {
		log.Fatalln("Not enough arguments!")
	}
	directory := args[0]
	checkEmpty(directory)
	Log.Verbose("init : creating new paranoid file system in " + directory)
	makeDir(directory, "names")
	makeDir(directory, "inodes")
	metaDir := makeDir(directory, "meta")
	makeDir(metaDir, "logs")
	makeDir(directory, "contents")
	uuid, err := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	uuidString := strings.TrimSpace(string(uuid))
	Log.Verbose("init uuid : " + uuidString)
	checkErr("init", err)
	err = ioutil.WriteFile(path.Join(metaDir, "uuid"), []byte(uuidString), 0600)
	checkErr("init", err)
	_, err = os.Create(path.Join(metaDir, "lock"))
	checkErr("init", err)
	io.WriteString(os.Stdout, returncodes.GetReturnCode(returncodes.OK))
}
