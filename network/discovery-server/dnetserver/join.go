package dnetserver

import (
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"log"
)

// Join method for Discovery Server
func (s *DiscoveryServer) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	log.Printf("[E] Join not yet implemented.\n")
	return nil, nil
}
