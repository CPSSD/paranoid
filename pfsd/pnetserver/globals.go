// Package pnetserver implements the ParanoidNetwork gRPC server.
// globals.go contains data used by each gRPC handler in pnetserver.
package pnetserver

import (
	"fmt"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

func convertCodeToError(code int, path string) error {
	switch code {
	case returncodes.EACCES:
		log.Printf("INFO: Do not have permission to edit %s.\n", path)
		returnError := grpc.Errorf(codes.PermissionDenied,
			"do not have permission to edit %s",
			path)
		return returnError
	case returncodes.ENOENT:
		log.Printf("INFO: File %s does not exist.\n", path)
		returnError := grpc.Errorf(codes.NotFound,
			"file %s does not exist",
			path)
		return returnError
	case returncodes.EEXIST:
		log.Printf("INFO: File %s already exists.\n", path)
		returnError := grpc.Errorf(codes.AlreadyExists,
			"file %s already exists",
			path)
		return returnError
	case returncodes.EISDIR:
		log.Printf("INFO: %s is a directory.\n", path)
		returnError := grpc.Errorf(codes.InvalidArgument,
			"%s is a directory",
			path)
		return returnError
	case returncodes.EIO:
		log.Println("INFO: Unexpected input or output from command.")
		returnError := grpc.Errorf(codes.Unknown, "Unexpected input or output")
		return returnError
	case returncodes.ENOTDIR:
		log.Printf("INFO: %s is not a directory.\n", path)
		returnError := grpc.Errorf(codes.InvalidArgument,
			"%s is not a directory",
			path)
		return returnError
	case returncodes.ENOTEMPTY:
		log.Printf("INFO: %s is not empty.\n", path)
		returnError := grpc.Errorf(codes.FailedPrecondition,
			"%s is not empty",
			path)
		return returnError
	}
	return nil
}
