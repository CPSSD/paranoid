package pnetserver

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
)

func (s *ParanoidServer) Unlink(ctx context.Context, req *pb.UnlinkRequest) (*pb.EmptyMessage, error) {
	code, err := commands.UnlinkCommand(ParanoidDir, req.Path, false)
	if code != returncodes.OK {
		log.Printf("ERROR: Could not unlink file %s: %v.\n", req.Path, err)
		returnError := grpc.Errorf(codes.Internal, "could not unlink file %s: %v", req.Path, err)
		return &pb.EmptyMessage{}, returnError
	}

	return &pb.EmptyMessage{}, nil
}
