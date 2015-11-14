package pnetserver

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func (s *ParanoidServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.EmptyMessage, error) {
	node := globals.Node{IP: req.Ip, Port: req.Port}
	log.Println("INFO: Got Ping from Node:", node)
	globals.Nodes.Add(node)
	return &pb.EmptyMessage{}, nil
}
