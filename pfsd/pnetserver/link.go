package pnetserver

import (
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func (s *ParanoidServer) Link(ctx context.Context, req *pb.LinkRequest) (*pb.EmptyMessage, error) {
	log.Printf("ERROR: Link not yet implemented.")
	return nil, nil
}
