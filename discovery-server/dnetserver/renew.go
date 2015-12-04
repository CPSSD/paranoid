package dnetserver

import (
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

// Renew method for Discovery Server
func (s *DiscoveryServer) Renew(ctx context.Context, req *pb.JoinRequest) (*pb.EmptyMessage, error) {
	isInNodes := false
	isActiveNode := false

	for i, node := range Nodes {
		if node.Data == *req.Node {
			if node.Active {
				Nodes[i].ExpiryTime = time.Now().Add(RenewInterval)
				isActiveNode = true
				Log.Infof("Renew: Node %s renewed", req.Node)
			}
			isInNodes = true
			break
		}
	}

	var returnError error
	if !isInNodes {
		returnError = grpc.Errorf(codes.NotFound, "node was not found")
	}
	if !isActiveNode {
		returnError = grpc.Errorf(codes.DeadlineExceeded, "renewal time is past deadline")
	}
	return &pb.EmptyMessage{}, returnError // returns nil if no error
}
