package pnetserver

import (
	"encoding/base64"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"strconv"
)

func (s *ParanoidServer) Write(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error) {
	var data []byte
	length, err := base64.StdEncoding.Decode(data, req.Data)
	if err != nil {
		log.Println("WARNING: Could not decode base64 data:", err)
	}
	var offset string
	// Go assumes a nil int field in a protobuf is 0
	if req.Offset != 0 {
		offset = strconv.FormatUint(req.Offset, 10)
	} else {
		offset = ""
	}
	code, _, err := runCommand(data, "write", ParanoidDir, req.Path, strconv.Itoa(length), offset)
	if err != nil {
		log.Printf("ERROR: Could not write to file %s: %v.\n", req.Path, err)
		returnError := grpc.Errorf(codes.Internal, "could not write to file %s: %v",
			req.Path, err)
		return &pb.WriteResponse{}, returnError
	}

	returnError := convertCodeToError(code, req.Path)
	// If returnError is nil here, it's equivalent to returning OK
	return &pb.WriteResponse{BytesWritten: uint64(length)}, returnError
}
