package pnetserver

import (
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func (s *ParanoidServer) Unlink(ctx context.Context, req *pb.UnlinkRequest) (*pb.EmptyMessage, error) {
	log.Printf("ERROR: Unlink not yet implemented.")
	return nil, nil
}
