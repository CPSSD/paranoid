package pnetserver

import (
	"encoding/base64"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
)

func (s *ParanoidServer) Utimes(ctx context.Context, req *pb.UtimesRequest) (*pb.EmptyMessage, error) {
	data := make([]byte, base64.StdEncoding.DecodedLen(len(req.Data)))
	realDecodedLen, err := base64.StdEncoding.Decode(data, req.Data)
	data = data[:realDecodedLen]
	if err != nil {
		log.Println("WARNING: Could not decode base64 data:", err)
	}

	code, _, err := runCommand(data, "utimes", ParanoidDir, req.Path)
	if err != nil {
		log.Printf("ERROR: Could not modify times of file %s: %v.\n", req.Path, err)
		returnError := grpc.Errorf(codes.Internal, "could not modify times of file %s: %v",
			req.Path, err)
		return &pb.EmptyMessage{}, returnError
	}

	returnError := convertCodeToError(code, req.Path)
	// If returnError is nil here, it's equivalent to returning OK
	return &pb.EmptyMessage{}, returnError
}
