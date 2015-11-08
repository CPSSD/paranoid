// Package pnetserver implements the ParanoidNetwork gRPC server.
// globals.go contains data used by each gRPC handler in pnetserver.
package pnetserver

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

type ParanoidServer struct{}

// Path to the PFS root directory
var ParanoidDir string

// RunCommand runs a pfs command with the given arguments. Gives stdinData on stdIn to pfs if it is not nil.
// Modified from github.com/cpssd/paranoid/pfi/pfsminterface.
func runCommand(stdinData []byte, cmdArgs ...string) (int, []byte, error) {
	cmdArgs = append(cmdArgs, "--net")
	command := exec.Command("pfsm", cmdArgs...)
	command.Stderr = os.Stderr

	if stdinData != nil {
		stdinPipe, err := command.StdinPipe()
		if err != nil {
			log.Println("ERROR: Could not open StdinPipe:", err)
			return -1, nil, fmt.Errorf("could not open StdinPipe:", err)
		}
		_, err = stdinPipe.Write(stdinData)
		if err != nil {
			log.Println("ERROR: Could not write to stdin:", err)
			return -1, nil, fmt.Errorf("could not write to stdin:", err)
		}
		err = stdinPipe.Close()
		if err != nil {
			log.Println("ERROR: Could not close StdinPipe:", err)
			return -1, nil, fmt.Errorf("could not close StdinPipe:", err)
		}
	}

	output, err := command.Output()
	if err != nil {
		log.Println("ERROR: Could not retrieve PFSM output:", err)
		return -1, nil, fmt.Errorf("could not retrieve PFSM output:", err)
	}
	code, err := strconv.Atoi(string(output[0:2]))
	if err != nil {
		log.Println("ERROR: PFSM returned invalid code:", err)
		return -1, nil, fmt.Errorf("PFSM returned invalid code:", err)
	}
	return code, output[2:], nil
}
