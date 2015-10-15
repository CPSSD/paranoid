package pfsInterface

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"
)

type statInfo struct {
	Length int64     `json:"length"`
	Ctime  time.Time `json:"ctime"`
	Mtime  time.Time `json:"mtime"`
	Atime  time.Time `json:"atime"`
}

/*
Stat -

description :
    Called when the attributes of a file or directory are needed.

paramenters :
    mountDir - The root directory of the file system.
    pfsLocation - The path to the pfs executable.
    name - The name of the file whos attributes are needed.

return :
    info - The statInfo object containing details of the file.
*/
func Stat(mountDir string, pfsLocation, name string) (info statInfo) {
	args := fmt.Sprintf("-f stat %s %s", mountDir, name)
	command := exec.Command(pfsLocation, args)

	output, err := command.Output()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))
	info = statInfo{}
	err = json.Unmarshal(output, &info)

	if err != nil {
		log.Fatal(err)
	}

	return info
	// TODO: return return structure object
}
