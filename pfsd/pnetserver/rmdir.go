package pnetserver

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (s *ParanoidServer) Rmdir(ctx context.Context, req *pb.RmdirRequest) (*pb.EmptyMessage, error) {
	code, err := commands.RmdirCommand(ParanoidDir, req.Directory, false)
	if code != returncodes.OK {
		Log.Errorf("Could not remove directory: %v \n", req.Directory, err)
		returnError := grpc.Errorf(codes.Internal, "could not remove directory: %v \n",
			req.Directory, err)
		return &pb.EmptyMessage{}, returnError
	}

	return &pb.EmptyMessage{}, nil
}
