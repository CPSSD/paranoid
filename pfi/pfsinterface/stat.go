package pfsinterface

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

//Stat gets the attributes of a file or directory from pfs
func Stat(initDir, name string) (info statInfo, e error) {
	command := exec.Command("pfs", "-f", "stat", initDir, name)
	output, err := command.Output()

	if err != nil {
		return statInfo{}, err
	}

	info = statInfo{}
	err = json.Unmarshal(output, &info)

	if err != nil {
		log.Fatalln(err)
	}

	return info, nil
}
