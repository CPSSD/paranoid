package pfsInterface

import (
	"log"
	"os/exec"
)

/*
Creat -

description :
    Called when a file is to be created.

parameters :
    initDir - The root directory of the pvd.
    pfsLocation - The path to the pfs executable.
    name - The name of the file to create.
*/
func Creat(initDir string, pfsLocation string, name string) {
	command := exec.Command(pfsLocation, "-f", "creat", initDir, name)

	err := command.Run()

	if err != nil {
		log.Fatal(err)
	}
}
