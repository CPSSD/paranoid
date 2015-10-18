package pfsinterface

import (
	"log"
	"os/exec"
	"strconv"
)

//Write tells pfs to write to a file
func Write(initDir string, pfsLocation string, name string, data []byte, offset, length int64) {
	var command *exec.Cmd
	if offset != -1 {
		command = exec.Command(pfsLocation, "-f", "write", initDir, name, strconv.FormatInt(offset, 10), strconv.FormatInt(length, 10))
	} else {
		command = exec.Command(pfsLocation, "-f", "write", initDir, name)
	}

	stdinPipe, err := command.StdinPipe()

	if err != nil {
		log.Fatalln(err)
	}

	err = command.Start()
	if err != nil {
		log.Fatalln(err)
	}

	stdinPipe.Write(data)
	stdinPipe.Close()
}
