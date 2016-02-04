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

	for i := 0; i < len(Nodes); i++ {
		if Nodes[i].Data.Uuid == req.Node.Uuid {
			if Nodes[i].Active {
				Nodes[i].ExpiryTime = time.Now().Add(RenewInterval)
				isActiveNode = true
				Log.Infof("Renew: Node %s (%s:%s) renewed", req.Node.Uuid, req.Node.Ip, req.Node.Port)
			}
			isInNodes = true
			break
		}
	}

	var returnError error
	if !isInNodes {
		Log.Errorf("Renew: node %s (%s:%s) not found", req.Node.Uuid, req.Node.Ip, req.Node.Port)
		returnError = grpc.Errorf(codes.NotFound, "node was not found")
	}
	if !isActiveNode {
		Log.Errorf("Renew: Renewal time is past deadline for node %s (%s:%s)", req.Node.Uuid, req.Node.Ip, req.Node.Port)
		returnError = grpc.Errorf(codes.DeadlineExceeded, "renewal time is past deadline")
	}
	return &pb.EmptyMessage{}, returnError // returns nil if no error
}
