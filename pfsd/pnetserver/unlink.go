package pnetserver

import (
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
)

func (s *ParanoidServer) Unlink(ctx context.Context, req *pb.UnlinkRequest) (*pb.EmptyMessage, error) {
	code, _, err := runCommand(nil, "unlink", ParanoidDir, req.Path)
	if err != nil {
		log.Printf("ERROR: Could not unlink file %s: %v.\n", req.Path, err)
		returnError := grpc.Errorf(codes.Internal, "could not unlink file %s: %v", req.Path, err)
		return &pb.EmptyMessage{}, returnError
	}

	returnError := convertCodeToError(code, req.Path)
	// If returnError is nil here, it's equivalent to returning OK
	return &pb.EmptyMessage{}, returnError
}
