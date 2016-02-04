package pnetserver

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"os"
)

func (s *ParanoidServer) Chmod(ctx context.Context, req *pb.ChmodRequest) (*pb.EmptyMessage, error) {
	code, err := commands.ChmodCommand(ParanoidDir, req.Path, os.FileMode(req.Mode), false)
	if code != returncodes.OK {
		Log.Errorf("Could not change permissions on file %s: %v.\n", req.Path, err)
		returnError := grpc.Errorf(codes.Internal, "could not change permissions on file %s: %v",
			req.Path, err)
		return &pb.EmptyMessage{}, returnError
	}

	return &pb.EmptyMessage{}, nil
}
