package pnetserver

import (
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func (s *ParanoidServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.EmptyMessage, error) {
	log.Printf("ERROR: Ping not yet implemented.")
	return nil, nil
}
