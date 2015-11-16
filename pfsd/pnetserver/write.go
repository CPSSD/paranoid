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
	_, err := base64.StdEncoding.Decode(data, req.Data)
	if err != nil {
		log.Println("WARNING: Could not decode base64 data:", err)
	}

	var code int
	var output []byte
	if req.Length != 0 {
		code, output, err = runCommand(data, "write", ParanoidDir, req.Path, strconv.FormatUint(req.Offset, 10), strconv.FormatUint(req.Length, 10))
	} else if req.Offset != 0 {
		code, output, err = runCommand(data, "write", ParanoidDir, req.Path, strconv.FormatUint(req.Offset, 10))
	} else {
		code, output, err = runCommand(data, "write", ParanoidDir, req.Path)
	}

	if err != nil {
		log.Printf("ERROR: Could not write to file %s: %v.\n", req.Path, err)
		returnError := grpc.Errorf(codes.Internal, "could not write to file %s: %v",
			req.Path, err)
		return &pb.WriteResponse{}, returnError
	}

	actualLength, err := strconv.ParseUint(string(output), 10, 64)
	if err != nil {
		log.Println("ERROR: Failed to convert length output:", err)
	}
	returnError := convertCodeToError(code, req.Path)
	// If returnError is nil here, it's equivalent to returning OK
	return &pb.WriteResponse{BytesWritten: actualLength}, returnError
}
