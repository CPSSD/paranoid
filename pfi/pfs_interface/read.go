package pfsInterface

import (
	"fmt"
	"log"
	"os/exec"
)

/*
Read -

description :
    Called when the contents of a file must be read.

parameters :
    initDir - The root directory of the pvd.
    pfsLocation - The path to the pfs executable.
    name - The name of the file to be read.
    offset (optional indicated by being -1) - The offset in bytes.
    length (optional indicated by being -1) - The length in bytes.

return :
    bytes - An array of bytes.
*/
func Read(initDir string, pfsLocation string, name string, offset int64, length int64) (bytes []byte) {
	command := exec.Command(pfsLocation, "-f", "read", initDir, name)
	if offset != -1 {
		command = exec.Command(pfsLocation, "-f", "read", initDir, name, fmt.Sprintf("%d", offset))
		if length != -1 {
			command = exec.Command(pfsLocation, "-f", "read", initDir, name, fmt.Sprintf("%d", offset), fmt.Sprintf("%d", length))
		}
	}

	output, err := command.Output()

	if err != nil {
		log.Fatal(err)
	}

	return output
}
