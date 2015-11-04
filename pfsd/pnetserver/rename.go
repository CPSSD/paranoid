package pnetserver

import (
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func (s *ParanoidServer) Rename(ctx context.Context, req *pb.RenameRequest) (*pb.EmptyMessage, error) {
	log.Printf("ERROR: Rename not yet implemented.")
	return nil, nil
}
