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
    initDir - The root directory of the pvd.
    pfsLocation - The path to the pfs executable.
    name - The name of the directory whose contents are needed.

return :
    fileNames - An array of strings representing the file names in the directory.
*/
func Readdir(initDir string, pfsLocation string, name string) (fileNames []string) {
	command := exec.Command(pfsLocation, "-f", "readdir", initDir)

	output, err := command.Output()
	outputString := string(output)

	if err != nil {
		log.Fatal(err)
	}
	if outputString == "" {
		return make([]string, 0)
	}

	filenames := strings.Split(outputString, "\n")
	fmt.Println(filenames)
	return filenames
}
