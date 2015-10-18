package pfsinterface

import (
	"log"
	"os/exec"
	"strconv"
)

//Read gets an array of bytes from pfs
func Read(initDir, pfsLocation, name string, offset, length int64) (bytes []byte) {
	var command *exec.Cmd
	if offset != -1 {
		command = exec.Command(pfsLocation, "-f", "read", initDir, name, strconv.FormatInt(offset, 10), strconv.FormatInt(length, 10))
	} else {
		command = exec.Command(pfsLocation, "-f", "read", initDir, name)
	}
	output, err := command.Output()

	if err != nil {
		log.Fatalln(err)
	}

	return output
}
