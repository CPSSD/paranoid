package pnetserver

import (
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func (s *ParanoidServer) Truncate(ctx context.Context, req *pb.TruncateRequest) (*pb.EmptyMessage, error) {
	log.Printf("ERROR: Truncate not yet implemented.")
	return nil, nil
}
