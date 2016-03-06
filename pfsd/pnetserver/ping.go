package pnetserver

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
)

func (s *ParanoidServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.EmptyMessage, error) {
	node := globals.Node{
		IP:         req.Ip,
		Port:       req.Port,
		CommonName: req.CommonName,
		UUID:       req.Uuid,
	}
	Log.Infof("Got Ping from Node:", node)
	globals.Nodes.Add(node)
	globals.RaftNetworkServer.ChangeNodeLocation(req.Uuid, req.Ip, req.Port)
	return &pb.EmptyMessage{}, nil
}
