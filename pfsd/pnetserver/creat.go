package pnetserver

import (
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"os/exec"
)

func (s *ParanoidServer) Creat(ctx context.Context, req *pb.CreatRequest) (*pb.EmptyMessage, error) {
	command := exec.Command("pfsm", "-n", "creat", PFSDir, req.Path)
	err := command.Run()
	if err != nil {
		log.Printf("ERROR: Could not create file %s: %v.\n", req.Path, err)
		returnError := grpc.Errorf(codes.Internal, "could not create file %s: %v", req.Path, err)
		return &pb.EmptyMessage{}, returnError
	}
	return &pb.EmptyMessage{}, nil
}
