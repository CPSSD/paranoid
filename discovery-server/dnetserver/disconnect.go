package dnetserver

import (
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Disconnect method for Discovery Server
func (s *DiscoveryServer) Disconnect(ctx context.Context, req *pb.DisconnectRequest) (*pb.EmptyMessage, error) {
	for i := 0; i < len(Nodes); i++ {
		if Nodes[i].Data.Uuid == req.Node.Uuid {
			Nodes[i].Active = false
			Log.Info("Disconnect: Node %s (%s:%s) disconnected", req.Node.Uuid, req.Node.Ip, req.Node.Port)
			return &pb.EmptyMessage{}, nil
		}
	}

	Log.Errorf("Disconnect: Node %s (%s:%s) was not found", req.Node.Uuid, req.Node.Ip, req.Node.Port)
	returnError := grpc.Errorf(codes.NotFound, "node was not found")
	return &pb.EmptyMessage{}, returnError
}
