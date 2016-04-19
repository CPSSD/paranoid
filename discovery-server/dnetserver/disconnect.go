package dnetserver

import (
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Disconnect method for Discovery Server
func (s *DiscoveryServer) Disconnect(ctx context.Context, req *pb.DisconnectRequest) (*pb.EmptyMessage, error) {
	PoolLock.RLock()
	defer PoolLock.RUnlock()

	if _, ok := Pools[req.Pool]; ok {
		Pools[req.Pool].PoolLock.Lock()
		defer Pools[req.Pool].PoolLock.Unlock()
		err := checkPoolPassword(req.Pool, req.Password, req.Node)
		if err != nil {
			return &pb.EmptyMessage{}, err
		}
	} else {
		Log.Errorf("Disconnect: Node %s (%s:%s) pool %s was not found", req.Node.Uuid, req.Node.Ip, req.Node.Port, req.Pool)
		returnError := grpc.Errorf(codes.NotFound, "pool %s was not found", req.Pool)
		return &pb.EmptyMessage{}, returnError
	}

	if _, ok := Pools[req.Pool].Info.Nodes[req.Node.Uuid]; ok {
		delete(Pools[req.Pool].Info.Nodes, req.Node.Uuid)
		saveState(req.Pool)
		Log.Info("Disconnect: Node %s (%s:%s) disconnected", req.Node.Uuid, req.Node.Ip, req.Node.Port)
		return &pb.EmptyMessage{}, nil
	}

	Log.Errorf("Disconnect: Node %s (%s:%s) was not found", req.Node.Uuid, req.Node.Ip, req.Node.Port)
	returnError := grpc.Errorf(codes.NotFound, "node was not found")
	return &pb.EmptyMessage{}, returnError
}
