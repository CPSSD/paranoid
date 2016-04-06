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
	defer saveState()
	nodes := getNodes(req.Pool, req.Node.Uuid)
	response := pb.JoinResponse{RenewInterval.Nanoseconds() / 1000 / 1000, nodes}

	seenPool, err := checkPoolPassword(req.Pool, req.Password, req.Node)
	if err != nil {
		return &pb.JoinResponse{}, err
	}

	if seenPool == false {
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
		newPool := Pool{
			Name:         req.Pool,
			PasswordSalt: salt,
			PasswordHash: hash,
		}
		Pools = append(Pools, newPool)
	}

	// Go through each node and check was the node there
	for i := 0; i < len(Nodes); i++ {
		if Nodes[i].Data.Uuid == req.Node.Uuid {
			if Nodes[i].Pool != req.Pool {
				Log.Errorf("Join: node belongs to pool %s but tried to join pool %s\n", Nodes[i].Pool, req.Pool)
				returnError := grpc.Errorf(codes.Internal,
					"node belongs to pool %s, but tried to join pool %s",
					Nodes[i].Pool, req.Pool)
				return &pb.JoinResponse{}, returnError
			}

			Nodes[i].Data.Ip = req.Node.Ip
			Nodes[i].Data.Port = req.Node.Port
			return &response, nil
		}
	}

	newNode := Node{req.Pool, *req.Node}
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
