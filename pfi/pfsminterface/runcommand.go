package pfsminterface

import (
	"log"
	"os"
	"os/exec"
	"strconv"
)

var OriginFlag string

//Current pfsm supported return codes
const (
	OK     = iota
	ENOENT //No such file or directory.
	EACCES //Can not access file
)

//RunCommand runs a pfs command with the given arguments. Gives stdinData on stdIn to pfs if it is not nil.
func RunCommand(stdinData []byte, cmdArgs ...string) (int, []byte) {
	cmdArgs = append(cmdArgs, OriginFlag)
	command := exec.Command("pfsm", cmdArgs...)
	command.Stderr = os.Stderr

	if stdinData != nil {
		stdinPipe, err := command.StdinPipe()
		if err != nil {
			log.Fatalln("Error running pfsm command :", err)
		}
		_, err = stdinPipe.Write(stdinData)
		if err != nil {
			log.Fatalln("Error running pfsm command :", err)
		}
		err = stdinPipe.Close()
		if err != nil {
			log.Fatalln("Error running pfsm command :", err)
		}
	}

	output, err := command.Output()
	if err != nil {
		log.Fatalln("Error running pfsm command :", err)
	}
	code, err := strconv.Atoi(string(output[0:2]))
	if err != nil {
		log.Fatalln("Invalid pfsm return code :", err)
	}
	return code, output[2:]
}
