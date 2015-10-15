package pfsInterface

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

/*
Read -

description :
    Called when the contents of a file must be read.

parameters :
    mountDir - The root directory of the file system.
    pfsLocation - The path to the pfs executable.
    name - The name of the file to be read.
    offset (optional indicated by being -1) - The offset in bytes.
    length (optional indicated by being -1) - The length in bytes.

return :
    bytes - An array of bytes.
*/
func Read(mountDir string, pfsLocation string, name string, offset int, length int) (bytes []byte) {
	args := fmt.Sprintf("-f read %s %s", mountDir, name)
	if offset != -1 {
		args += " " + string(offset)

		if length != -1 {
			args += " " + string(length)
		}
	}

	command := exec.Command(pfsLocation, args)

	output, err := command.Output()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println(output)
	return output
}
