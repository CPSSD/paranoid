package pnetserver

import (
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
)

func (s *ParanoidServer) Rmdir(ctx context.Context, req *pb.RmdirRequest) (*pb.EmptyMessage, error) {
	code, _, err := runCommand(nil, "rmdir", ParanoidDir, req.directory)
	if err != nil {
		log.Printf("ERROR: Could not remove directory: %v \n", req.directory, err)
		returnError := grpc.Errorf(codes.Internal, "Could not remove directory: %v \n",
			req.directory, err)
		return &pb.EmptyMessage{}, returnError
	}

	returnError := convertCodeToError(code, req.directory)
	// If returnError is nil here, it's equivalent to returning OK
	return &pb.EmptyMessage{}, returnError
}
