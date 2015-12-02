package pnetserver

import (
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"strconv"
)

func (s *ParanoidServer) Mkdir(ctx context.Context, req *pb.MkdirRequest) (*pb.EmptyMessage, error) {
	code, _, err := runCommand(nil, "mkdir", ParanoidDir, req.Directory, strconv.FormatUint(uint64(req.Mode), 8))
	if err != nil {
		log.Printf("ERROR: Could not make directory: %v with mode: %v \n", req.Directory, req.Mode, err)
		returnError := grpc.Errorf(codes.Internal, "Could not make directory: %v with mode: %v\n",
			req.Directory, req.Mode, err)
		return &pb.EmptyMessage{}, returnError
	}

	returnError := convertCodeToError(code, req.Directory)
	// If returnError is nil here, it's equivalent to returning OK
	return &pb.EmptyMessage{}, returnError
}
