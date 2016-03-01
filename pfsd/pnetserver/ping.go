package pnetserver

import (
	"fmt"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"github.com/cpssd/paranoid/raft"
	"golang.org/x/net/context"
)

//Recieve a ping from a node asking to join raft network
func (s *ParanoidServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.EmptyMessage, error) {
	node := raft.Node{
		IP:         req.Ip,
		Port:       req.Port,
		CommonName: req.CommonName,
		NodeID:     req.Uuid,
	}
	Log.Infof("Got Ping from Node:", node)
	err := RaftNetworkServer.RequestAddNodeToConfiguration(node)
	if err != nil {
		return &pb.EmptyMessage{}, fmt.Errorf("unable to add node to raft cluster:", err)
	}
	return &pb.EmptyMessage{}, nil
}
