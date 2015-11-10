package pnetserver

import (
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
)

func (s *ParanoidServer) Link(ctx context.Context, req *pb.LinkRequest) (*pb.EmptyMessage, error) {
	code, _, err := runCommand(nil, "link", ParanoidDir, req.OldPath, req.NewPath)
	if err != nil {
		log.Printf("ERROR: Could not link file %s to %s: %v.\n", req.OldPath, req.NewPath, err)
		returnError := grpc.Errorf(codes.Internal, "could not link file %s to %s: %v",
			req.OldPath, req.NewPath, err)
		return &pb.EmptyMessage{}, returnError
	}

	returnError := convertCodeToError(code, req.OldPath)
	// If returnError is nil here, it's equivalent to returning OK
	return &pb.EmptyMessage{}, returnError
}
