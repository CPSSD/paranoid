package pfsInterface

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

/*
Readdir -

description :
    Called when the contents of a directory are needed.

parameters :
    mountDir - The root directory of the file system.
    pfsLocation - The path to the pfs executable.
    name - The name of the directory whose contents are needed.

return :
    fileNames - An array of strings representing the file names in the directory.
*/
func Readdir(mountDir string, pfsLocation string, name string) (fileNames []string) {
	args := fmt.Sprintf("-f readdir %s", mountDir)
	command := exec.Command(pfsLocation, args)

	output, err := command.Output()

	if err != nil {
		log.Fatal(err)
	}

	outputString := string(output)
	filenames := strings.Split(outputString, "\n")
	fmt.Println(filenames)
	return filenames
}
