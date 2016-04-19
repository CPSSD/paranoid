package dnetserver

import (
	"crypto/rand"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"io"
)

const (
	PASSWORD_SALT_LENGTH int = 64
)

// Join method for Discovery Server
func (s *DiscoveryServer) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	PoolLock.RLock()
	if _, ok := Pools[req.Pool]; ok {
		defer PoolLock.RUnlock()
		Pools[req.Pool].PoolLock.Lock()
		defer Pools[req.Pool].PoolLock.Unlock()

		err := checkPoolPassword(req.Pool, req.Password, req.Node)
		if err != nil {
			return &pb.JoinResponse{}, err
		}
	} else {
		PoolLock.RUnlock()
		PoolLock.Lock()
		defer PoolLock.Unlock()

		if _, ok := Pools[req.Pool]; ok {
			Pools[req.Pool].PoolLock.Lock()
			defer Pools[req.Pool].PoolLock.Unlock()
			err := checkPoolPassword(req.Pool, req.Password, req.Node)
			if err != nil {
				return &pb.JoinResponse{}, err
			}
		} else {
			hash := make([]byte, 0)
			salt := make([]byte, PASSWORD_SALT_LENGTH)
			n, err := io.ReadFull(rand.Reader, salt)
			if err != nil {
				returnError := grpc.Errorf(codes.Internal,
					"error hashing password: %s",
					err,
				)
				return &pb.JoinResponse{}, returnError
			}
			if n != PASSWORD_SALT_LENGTH {
				returnError := grpc.Errorf(codes.Internal,
					"error hashing password: unable to read salt from random number generator",
				)
				return &pb.JoinResponse{}, returnError
			}

			if req.Password != "" {
				hash, err = bcrypt.GenerateFromPassword(append(salt, []byte(req.Password)...), bcrypt.DefaultCost)
				if err != nil {
					returnError := grpc.Errorf(codes.Internal,
						"error hashing password: %s",
						err,
					)
					return &pb.JoinResponse{}, returnError
				}
			}
			newPool := &Pool{
				Info: PoolInfo{
					PasswordSalt: salt,
					PasswordHash: hash,
				},
			}
			Pools[req.Pool] = newPool
			Pools[req.Pool].Info.Nodes = make(map[string]*pb.Node)
			Pools[req.Pool].PoolLock.Lock()
			defer Pools[req.Pool].PoolLock.Unlock()
		}
	}

	nodes := getNodes(req.Pool, req.Node.Uuid)
	response := pb.JoinResponse{RenewInterval.Nanoseconds() / 1000 / 1000, nodes}

	Pools[req.Pool].Info.Nodes[req.Node.Uuid] = req.Node
	Log.Infof("Join: Node %s (%s:%s) joined \n", req.Node.Uuid, req.Node.Ip, req.Node.Port)
	saveState(req.Pool)

	return &response, nil
}

func getNodes(pool, requesterUuid string) []*pb.Node {
	var nodes []*pb.Node
	if _, ok := Pools[pool]; ok {
		for nodeUUID, _ := range Pools[pool].Info.Nodes {
			nodes = append(nodes, Pools[pool].Info.Nodes[nodeUUID])
		}
	}
	return nodes
}
