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

func (s *ParanoidServer) Write(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error) {
	code, err, bytesWritten := commands.WriteCommand(ParanoidDir, req.Path, int64(req.Offset), int64(req.Length), req.Data, false)

	if code != returncodes.OK {
		log.Printf("ERROR: Could not write to file %s: %v.\n", req.Path, err)
		returnError := grpc.Errorf(codes.Internal, "could not write to file %s: %v",
			req.Path, err)
		return &pb.WriteResponse{}, returnError
	}

	return &pb.WriteResponse{BytesWritten: uint64(bytesWritten)}, nil
}
