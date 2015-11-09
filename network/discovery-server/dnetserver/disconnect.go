package dnetserver

import (
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"log"
)

// Disconnect method for Discovery Server
func (s *DiscoveryServer) Disconnect(ctx context.Context, req *pb.DisconnectRequest) (*pb.EmptyMessage, error) {
	log.Printf("[E] Disconnect not yet implemented.\n")
	return nil, nil
}
