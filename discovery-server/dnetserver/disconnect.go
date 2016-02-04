package dnetserver

import (
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Disconnect method for Discovery Server
func (s *DiscoveryServer) Disconnect(ctx context.Context, req *pb.DisconnectRequest) (*pb.EmptyMessage, error) {
	for i, node := range Nodes {
		if node.Data == *req.Node {
			Nodes[i].Active = false
			Log.Info("Disconnect: Node %s:%s disconnected", req.Node.Ip, req.Node.Port)
			return &pb.EmptyMessage{}, nil
		}
	}

	Log.Errorf("Disconnect: Node %s:%s was not found", req.Node.Ip, req.Node.Port)
	returnError := grpc.Errorf(codes.NotFound, "node was not found")
	return &pb.EmptyMessage{}, returnError
}
