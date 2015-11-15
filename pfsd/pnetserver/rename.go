package pnetserver

import (
	"github.com/cpssd/paranoid/pfsm/returncodes"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
)

func (s *ParanoidServer) Rename(ctx context.Context, req *pb.RenameRequest) (*pb.EmptyMessage, error) {
	code, _, err := runCommand(nil, "rename", ParanoidDir, req.OldPath, req.NewPath)
	if err != nil {
		log.Printf("ERROR: Could not rename file %s: %v.\n", req.OldPath, err)
		returnError := grpc.Errorf(codes.Internal, "could not rename file %s: %v",
			req.OldPath, err)
		return &pb.EmptyMessage{}, returnError
	}

	if code == returncodes.EEXIST {
		log.Printf("INFO: File %s already exists.\n", req.NewPath)
		returnError := grpc.Errorf(codes.AlreadyExists,
			"file %s already exists",
			req.NewPath)
		return &pb.EmptyMessage{}, returnError
	}
	returnError := convertCodeToError(code, req.OldPath)
	// If returnError is nil here, it's equivalent to returning OK
	return &pb.EmptyMessage{}, returnError
}
