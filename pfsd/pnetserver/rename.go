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

func (s *ParanoidServer) Rename(ctx context.Context, req *pb.RenameRequest) (*pb.EmptyMessage, error) {
	code, err := commands.RenameCommand(ParanoidDir, req.OldPath, req.NewPath, false)
	if code != returncodes.OK {
		log.Printf("ERROR: Could not rename file %s: %v.\n", req.OldPath, err)
		returnError := grpc.Errorf(codes.Internal, "could not rename file %s: %v",
			req.OldPath, err)
		return &pb.EmptyMessage{}, returnError
	}

	return &pb.EmptyMessage{}, nil
}
