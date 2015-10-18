package pfsinterface

import (
	"log"
	"os/exec"
)

//Creat tells pfs that a file needs to be created
func Creat(initDir, pfsLocation, name string) {
	command := exec.Command(pfsLocation, "-f", "creat", initDir, name)

	err := command.Run()

	if err != nil {
		log.Fatalln(err)
	}
}
