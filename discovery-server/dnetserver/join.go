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
	nodes := getNodes(req.Pool)
	response := pb.JoinResponse{RenewInterval.Nanoseconds() / 1000 / 1000, nodes}

	// Go through each node and check was the node there
	for _, node := range Nodes {
		if node.Data == *req.Node {
			if node.Active {
				Log.Errorf("Join: node %s:%s is already part of the cluster", req.Node.Ip, req.Node.Port)
				returnError := grpc.Errorf(codes.AlreadyExists,
					"node is already part of the cluster")
				return &pb.JoinResponse{}, returnError
			}

			if node.Pool != req.Pool {
				Log.Errorf("Join: node belongs to pool %s but tried to join pool %s\n", node.Pool, req.Pool)
				returnError := grpc.Errorf(codes.Internal,
					"node belongs to pool %s, but tried to join pool %s",
					node.Pool, req.Pool)
				return &pb.JoinResponse{}, returnError
			}

			node.Active = true
			return &response, nil
		}
	}

	newNode := Node{true, req.Pool, time.Now().Add(RenewInterval), *req.Node}
	Nodes = append(Nodes, newNode)
	Log.Infof("Join: Node %s:%s joined \n", req.Node.Ip, req.Node.Port)

	return &response, nil
}

func getNodes(pool string) []*pb.Node {
	var nodes []*pb.Node
	for i := 0; i < len(Nodes); i++ {
		if Nodes[i].Active && Nodes[i].Pool == pool {
			nodes = append(nodes, &(Nodes[i].Data))
		}
	}
	return nodes
}
