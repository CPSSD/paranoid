package pfsminterface

import (
	"log"
	"os"
	"os/exec"
)

var OriginFlag string

//RunCommand runs a pfs command with the given arguments. Gives stdinData on stdIn to pfs if it is not nil.
func RunCommand(stdinData []byte, cmdArgs ...string) []byte {
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
	return output
}
