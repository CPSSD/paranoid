package pfsInterface

import (
	"fmt"
	"log"
	"os/exec"
)

/*
Creat -

description :
    Called when a file is to be created.

parameters :
    mountDir - The root directory of the file system.
    pfsLocation - The path to the pfs executable.
    name - The name of the file to create.
*/
func Creat(mountDir string, pfsLocation string, name string) {
	args := fmt.Sprintf("-f creat %s %s", mountDir, name)
	command := exec.Command(pfsLocation, args)

	err := command.Run()

	if err != nil {
		log.Fatal(err)
	}
}
