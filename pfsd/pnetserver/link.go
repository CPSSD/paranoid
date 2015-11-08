package pnetserver

import (
	"bufio"
	"github.com/cpssd/paranoid/pfsm/returncodes"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"os/exec"
	"strconv"
)

func (s *ParanoidServer) Link(ctx context.Context, req *pb.LinkRequest) (*pb.EmptyMessage, error) {
	command := exec.Command("pfsm", "-n", "link", ParanoidDir, req.OldPath, req.NewPath)
	stdout, err := command.StdoutPipe()
	if err != nil {
		log.Println("ERROR: Could not capture stdout of subprocess.")
	}
	read := bufio.NewReader(stdout)
	err = command.Run()
	if err != nil {
		log.Printf("ERROR: Could not link file %s to %s: %v.\n", req.OldPath, req.NewPath, err)
		returnError := grpc.Errorf(codes.Internal, "could not link file %s to %s: %v",
			req.OldPath, req.NewPath, err)
		return &pb.EmptyMessage{}, returnError
	}

	pfsmoutput, _, err := read.ReadLine()
	codeBytes := pfsmoutput[:2]
	code, err := strconv.Atoi(string(codeBytes))
	if err != nil {
		log.Println("ERROR: Could not interpret PFSM error code.")
		returnError := grpc.Errorf(codes.Internal, "could not link file %s to %s: %v",
			req.OldPath, req.NewPath, err)
		return &pb.EmptyMessage{}, returnError
	}

	switch code {
	case returncodes.EACCES:
		log.Printf("INFO: Do not have permission to edit %s.\n", req.OldPath)
		returnError := grpc.Errorf(codes.PermissionDenied,
			"do not have permission to edit %s",
			req.OldPath)
		return &pb.EmptyMessage{}, returnError
	case returncodes.ENOENT:
		log.Printf("INFO: File %s does not exist.\n", req.OldPath)
		returnError := grpc.Errorf(codes.NotFound,
			"file %s does not exist",
			req.OldPath)
		return &pb.EmptyMessage{}, returnError
	}

	return &pb.EmptyMessage{}, nil
}
