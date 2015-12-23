package pnetserver

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"os"
)

func (s *ParanoidServer) Mkdir(ctx context.Context, req *pb.MkdirRequest) (*pb.EmptyMessage, error) {
	code, err := commands.MkdirCommand(ParanoidDir, req.Directory, os.FileMode(req.Mode), false)
	if code != returncodes.OK {
		log.Printf("ERROR: Could not make directory: %v with mode: %v \n", req.Directory, req.Mode, err)
		returnError := grpc.Errorf(codes.Internal, "could not make directory: %v with mode: %v\n",
			req.Directory, req.Mode, err)
		return &pb.EmptyMessage{}, returnError
	}

	return &pb.EmptyMessage{}, nil
}
