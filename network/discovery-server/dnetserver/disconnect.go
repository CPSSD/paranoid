package dnetserver

import (
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"reflect"
)

// Disconnect method for Discovery Server
func (s *DiscoveryServer) Disconnect(ctx context.Context, req *pb.DisconnectRequest) (*pb.EmptyMessage, error) {
	isInNodes := false

	for i, node := range Nodes {
		if reflect.DeepEqual(&node, req.Node) {
			Nodes[i].Active = false
			isInNodes = true
		}
	}

	var returnError error
	if !isInNodes {
		returnError = grpc.Errorf(codes.NotFound, "node was not found")
	}
	return &pb.EmptyMessage{}, returnError // returns nil if no error
}
