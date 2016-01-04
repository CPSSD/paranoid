package pnetserver

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"time"
)

func (s *ParanoidServer) Utimes(ctx context.Context, req *pb.UtimesRequest) (*pb.EmptyMessage, error) {
	var atime *time.Time
	var mtime *time.Time
	if req.AccessNanoseconds != 0 || req.AccessSeconds != 0 {
		time := time.Unix(req.AccessSeconds, req.AccessNanoseconds)
		atime = &time
	}
	if req.ModifyNanoseconds != 0 || req.ModifySeconds != 0 {
		time := time.Unix(req.ModifySeconds, req.ModifyNanoseconds)
		mtime = &time
	}
	code, err := commands.UtimesCommand(ParanoidDir, req.Path, atime, mtime, false)
	if code != returncodes.OK {
		log.Printf("ERROR: Could not modify times of file %s: %v.\n", req.Path, err)
		returnError := grpc.Errorf(codes.Internal, "could not modify times of file %s: %v",
			req.Path, err)
		return &pb.EmptyMessage{}, returnError
	}

	return &pb.EmptyMessage{}, nil
}
