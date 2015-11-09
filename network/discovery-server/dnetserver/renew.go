package dnetserver

import (
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"log"
)

// Renew method for Discovery Server
func (s *DiscoveryServer) Renew(ctx context.Context, req *pb.JoinRequest) (*pb.EmptyMessage, error) {
	log.Printf("[E] Renew not yet implemented.\n")
	return nil, nil
}
