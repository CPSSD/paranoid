package pnetserver

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (s *ParanoidServer) Truncate(ctx context.Context, req *pb.TruncateRequest) (*pb.EmptyMessage, error) {
	code, err := commands.TruncateCommand(ParanoidDir, req.Path, int64(req.Length), false)
	if code != returncodes.OK {
		Log.Errorf("Could not truncate file %s: %v.\n", req.Path, err)
		returnError := grpc.Errorf(codes.Internal, "could not truncate file %s: %v",
			req.Path, err)
		return &pb.EmptyMessage{}, returnError
	}

	return &pb.EmptyMessage{}, nil
}
