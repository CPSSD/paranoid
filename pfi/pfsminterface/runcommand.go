package pfsminterface

import (
	"github.com/cpssd/paranoid/pfi/util"
	"os"
	"os/exec"
	"strconv"
)

var OriginFlag string

//RunCommand runs a pfs command with the given arguments. Gives stdinData on stdIn to pfs if it is not nil.
func RunCommand(stdinData []byte, cmdArgs ...string) (int, []byte) {
	cmdArgs = append(cmdArgs, OriginFlag)
	if util.LogOutput {
		cmdArgs = append(cmdArgs, "-v")
	}
	command := exec.Command("pfsm", cmdArgs...)
	command.Stderr = os.Stderr

	if stdinData != nil {
		stdinPipe, err := command.StdinPipe()
		if err != nil {
			util.Log.Fatal("Error running pfsm command :", err)
		}
		_, err = stdinPipe.Write(stdinData)
		if err != nil {
			util.Log.Fatal("Error running pfsm command :", err)
		}
		err = stdinPipe.Close()
		if err != nil {
			util.Log.Fatal("Error running pfsm command :", err)
		}
	}

	output, err := command.Output()
	if err != nil {
		util.Log.Fatal("Error running pfsm command :", err)
	}
	code, err := strconv.Atoi(string(output[0:2]))
	if err != nil {
		util.Log.Fatal("Invalid pfsm return code :", err)
	}
	return code, output[2:]
}
