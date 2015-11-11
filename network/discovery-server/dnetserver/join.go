package dnetserver

import (
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"reflect"
	"time"
)

// Join method for Discovery Server
func (s *DiscoveryServer) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	// Check was the node there
	check := 1

	// Go through each node and check was the node there
	for _, node := range Nodes {
		if reflect.DeepEqual(&node.Data, req.Node) {
			if node.Active {
				returnError := grpc.Errorf(codes.AlreadyExists,
					"node is already part of the cluster")
				return &pb.JoinResponse{}, returnError
			}

			if node.Pool != req.Pool {
				returnError := grpc.Errorf(codes.Internal,
					"node belongs to pool %s, but tried to join pool %s",
					node.Pool, req.Pool)
				return &pb.JoinResponse{}, returnError
			}

			node.Active = true
			check = 0
		}
	}

	nodes := getNodes(req.Pool)
	response := pb.JoinResponse{RenewInterval.Nanoseconds() * 1000 * 1000, nodes}

	if check == 1 {
		newNode := Node{
			true,
			req.Pool,
			time.Now().Add(RenewInterval),
			*req.Node}

		Nodes = append(Nodes, newNode)
	}

	return &response, nil
}

func getNodes(pool string) []*pb.Node {
	var nodes []*pb.Node
	for _, node := range Nodes {
		if node.Pool == pool {
			nodes = append(nodes, &node.Data)
		}
	}
	return nodes
}
