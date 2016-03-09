package pnetserver

import (
	"fmt"
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"github.com/cpssd/paranoid/raft"
	"golang.org/x/net/context"
)

//JoinCluster recieves requests from nodes asking to join raft cluster
func (s *ParanoidServer) JoinCluster(ctx context.Context, req *pb.PingRequest) (*pb.EmptyMessage, error) {
	node := raft.Node{
		IP:         req.Ip,
		Port:       req.Port,
		CommonName: req.CommonName,
		NodeID:     req.Uuid,
	}
	Log.Infof("Got Ping from Node:", node)
	err := globals.RaftNetworkServer.RequestAddNodeToConfiguration(node)
	if err != nil {
		return &pb.EmptyMessage{}, fmt.Errorf("unable to add node to raft cluster:", err)
	}
	return &pb.EmptyMessage{}, nil
}
