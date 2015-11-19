package pnetserver

import (
	"encoding/json"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"time"
)

type utimesTime struct {
	Atime time.Time `json:"atime",omitempty`
	Mtime time.Time `json:"mtime",omitempty`
}

func (s *ParanoidServer) Utimes(ctx context.Context, req *pb.UtimesRequest) (*pb.EmptyMessage, error) {
	timeStruct := &utimesTime{}
	if req.AccessSeconds != 0 || req.AccessMicroseconds != 0 {
		timeStruct.Atime = time.Unix(int64(req.AccessSeconds), int64(req.AccessMicroseconds)*1000)
	}
	if req.ModifySeconds != 0 || req.ModifyMicroseconds != 0 {
		timeStruct.Mtime = time.Unix(int64(req.ModifySeconds), int64(req.ModifyMicroseconds)*1000)
	}
	data, err := json.Marshal(timeStruct)
	if err != nil {
		log.Printf("WARNING: Error marshaling time to JSON:", err)
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
