package pnetserver

import (
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func (s *ParanoidServer) Write(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error) {
	log.Printf("ERROR: Write not yet implemented.")
	return nil, nil
}
