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
    mountDir - The root directory of the file system.
    pfsLocation - The path to the pfs executable.
    name - The name of the file to be read.
    data - The data (array of bytes) to write to the file.
    offset (optional indicated by being -1) - The offset in bytes.
    length (optional indicated by being -1) - The length in bytes.
*/
func Write(mountDir string, pfsLocation string, name string, data []byte, offset int, length int) {
	args := fmt.Sprintf("-f write %s %s", mountDir, name)
	if offset != -1 {
		args += " " + string(offset)

		if length != -1 {
			args += " " + string(length)
		}
	}

	command := exec.Command(pfsLocation, args)
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
