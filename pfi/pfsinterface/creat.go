package pfsinterface

import (
	"log"
	"os/exec"
)

var OriginFlag string

//Creat tells pfs that a file needs to be created
func Creat(initDir, name string) {
	command := exec.Command("pfs", OriginFlag, "creat", initDir, name)

	err := command.Run()

	if err != nil {
		log.Fatalln(err)
	}
}
