package pfsInterface

import (
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
func Write(initDir string, pfsLocation string, name string, data []byte, offset int, length int) {
	command := exec.Command(pfsLocation, "-f", "write", initDir, name)
	if offset != -1 {
		command = exec.Command(pfsLocation, "-f", "read", initDir, name, string(offset))
		if length != -1 {
			command = exec.Command(pfsLocation, "-f", "read", initDir, name, string(offset), string(length))
		}
	}

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
