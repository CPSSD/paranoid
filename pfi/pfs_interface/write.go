package pfsInterface

import (
	"fmt"
	"log"
	"os/exec"
)

/*
Write -

description :
    Called when writing to a file.

parameters :
    initDir - The root directory of the pvd.
    pfsLocation - The path to the pfs executable.
    name - The name of the file to be read.
    data - The data (array of bytes) to write to the file.
    offset (optional indicated by being -1) - The offset in bytes.
    length (optional indicated by being -1) - The length in bytes.
*/
func Write(initDir string, pfsLocation string, name string, data []byte, offset int64, length int64) {
	command := exec.Command(pfsLocation, "-f", "write", initDir, name)
	if offset != -1 {
		command = exec.Command(pfsLocation, "-f", "write", initDir, name, fmt.Sprintf("%d", offset))
		if length != -1 {
			command = exec.Command(pfsLocation, "-f", "write", initDir, name, fmt.Sprintf("%d", offset), fmt.Sprintf("%d", length))
		}
	}

	fmt.Println(pfsLocation, "-f", "write", initDir, name, fmt.Sprintf("%d", offset), fmt.Sprintf("%d", length))

	stdinPipe, err := command.StdinPipe()

	if err != nil {
		log.Fatal(err)
	}

	err = command.Start()
	if err != nil {
		log.Fatal(err)
	}

	stdinPipe.Write(data)
	stdinPipe.Close()
}
