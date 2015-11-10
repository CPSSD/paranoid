package pnetserver

import (
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
)

func (s *ParanoidServer) Truncate(ctx context.Context, req *pb.TruncateRequest) (*pb.EmptyMessage, error) {
	code, _, err := runCommand(nil, "truncate", ParanoidDir, req.Path, string(req.Length))
	if err != nil {
		log.Printf("ERROR: Could not truncate file %s: %v.\n", req.Path, err)
		returnError := grpc.Errorf(codes.Internal, "could not truncate file %s: %v",
			req.Path, err)
		return &pb.EmptyMessage{}, returnError
	}

	returnError := convertCodeToError(code, req.Path)
	// If returnError is nil here, it's equivalent to returning OK
	return &pb.EmptyMessage{}, returnError
}
