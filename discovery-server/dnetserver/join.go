package dnetserver

import (
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

// Join method for Discovery Server
func (s *DiscoveryServer) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	nodes := getNodes(req.Pool, req.Node.Uuid)
	response := pb.JoinResponse{RenewInterval.Nanoseconds() / 1000 / 1000, nodes}

	// Go through each node and check was the node there
	for i := 0; i < len(Nodes); i++ {
		if Nodes[i].Data.Uuid == req.Node.Uuid {
			if Nodes[i].Active {
				Log.Errorf("Join: node %s (%s:%s) is already part of the cluster", req.Node.Uuid, req.Node.Ip, req.Node.Port)
				returnError := grpc.Errorf(codes.AlreadyExists,
					"node is already part of the cluster")
				return &pb.JoinResponse{}, returnError
			}

			if Nodes[i].Pool != req.Pool {
				Log.Errorf("Join: node belongs to pool %s but tried to join pool %s\n", Nodes[i].Pool, req.Pool)
				returnError := grpc.Errorf(codes.Internal,
					"node belongs to pool %s, but tried to join pool %s",
					Nodes[i].Pool, req.Pool)
				return &pb.JoinResponse{}, returnError
			}

			Nodes[i].Active = true
			return &response, nil
		}
	}

	newNode := Node{true, req.Pool, time.Now().Add(RenewInterval), *req.Node}
	Nodes = append(Nodes, newNode)
	Log.Infof("Join: Node %s (%s:%s) joined \n", req.Node.Uuid, req.Node.Ip, req.Node.Port)

	return &response, nil
}

func getNodes(pool, requesterUuid string) []*pb.Node {
	var nodes []*pb.Node
	for i := 0; i < len(Nodes); i++ {
		if Nodes[i].Pool == pool && Nodes[i].Data.Uuid != requesterUuid {
			nodes = append(nodes, &(Nodes[i].Data))
		}
	}
	return nodes
}
