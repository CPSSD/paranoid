package pnetserver

import (
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func (s *ParanoidServer) Chmod(ctx context.Context, req *pb.ChmodRequest) (*pb.EmptyMessage, error) {
	log.Printf("ERROR: Chmod not yet implemented.")
	return nil, nil
}
