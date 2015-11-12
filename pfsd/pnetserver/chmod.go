package pnetserver

import (
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"strconv"
)

func (s *ParanoidServer) Chmod(ctx context.Context, req *pb.ChmodRequest) (*pb.EmptyMessage, error) {
	code, _, err := runCommand(nil, "chmod", ParanoidDir, req.Path, strconv.FormatUint(uint64(req.Mode), 8))
	if err != nil {
		log.Printf("ERROR: Could not change permissions on file %s: %v.\n", req.Path, err)
		returnError := grpc.Errorf(codes.Internal, "could not change permissions on file %s: %v",
			req.Path, err)
		return &pb.EmptyMessage{}, returnError
	}

	returnError := convertCodeToError(code, req.Path)
	// If returnError is nil here, it's equivalent to returning OK
	return &pb.EmptyMessage{}, returnError
}
