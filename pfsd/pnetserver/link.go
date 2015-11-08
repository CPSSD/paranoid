package pnetserver

import (
	"github.com/cpssd/paranoid/pfsm/returncodes"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
)

func (s *ParanoidServer) Link(ctx context.Context, req *pb.LinkRequest) (*pb.EmptyMessage, error) {
	code, _, err := runCommand(nil, "link", req.OldPath, req.NewPath)
	if err != nil {
		log.Printf("ERROR: Could not link file %s to %s: %v.\n", req.OldPath, req.NewPath, err)
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
