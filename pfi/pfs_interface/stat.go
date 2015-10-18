package pfsInterface

import (
	"encoding/json"
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
    initDir - The root directory of the pvd.
    pfsLocation - The path to the pfs executable.
    name - The name of the file whos attributes are needed.

return :
    info - The statInfo object containing details of the file.
*/
func Stat(initDir string, pfsLocation, name string) (info statInfo, e error) {
	command := exec.Command(pfsLocation, "-f", "stat", initDir, name)
	output, err := command.Output()

	if err != nil {
		return statInfo{}, err
	}

	info = statInfo{}
	err = json.Unmarshal(output, &info)

	if err != nil {
		log.Fatal(err)
	}

	return info, nil
}
