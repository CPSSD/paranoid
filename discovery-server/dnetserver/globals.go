// Package dnetserver implements the DiscoveryNetwork gRPC server.
// globals.go contains data used by each gRPC handler in dnetserver.
package dnetserver

import (
	"github.com/cpssd/paranoid/logger"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

var Log *logger.ParanoidLogger

// DiscoveryServer struct
type DiscoveryServer struct{}

// Node struct to hold the node data
type Node struct {
	Pool string  `json:"pool"`
	Data pb.Node `json:"data"`
}

// Pool struct to hold the pool data
type Pool struct {
	Name         string `json:"name"`
	PasswordSalt []byte `json:"passwordsalt"`
	PasswordHash []byte `json:"passwordhash"`
}

// Nodes array
var Nodes []Node

//Pools array
var Pools []Pool

// RenewInterval global containing the time after which the nodes will be marked as inactive
var RenewInterval time.Duration

// StateFilePath is the path to the file in which the discovery server stores its state
var StateFilePath string

func checkPoolPassword(pool, password string, node *pb.Node) (seenPool bool, err error) {
	seenPool = false
	for i := 0; i < len(Pools); i++ {
		if Pools[i].Name == pool {
			if password == "" {
				if len(Pools[i].PasswordHash) != 0 {
					Log.Errorf("Join: node %s attempted join password protected pool without a giving a password", node.Uuid)
					returnError := grpc.Errorf(codes.Internal,
						"pool %s is password protected",
						pool,
					)
					return true, returnError
				}
			} else {
				err := bcrypt.CompareHashAndPassword(Pools[i].PasswordHash, append(Pools[i].PasswordSalt, []byte(password)...))
				if err != nil {
					Log.Errorf("Join: node %s attempted join password protected pool with incorrect password: %s",
						node.Uuid,
						err,
					)
					returnError := grpc.Errorf(codes.Internal,
						"given password incorrect: %s",
						err,
					)
					return true, returnError
				}
			}
			seenPool = true
			break
		}
	}
	return seenPool, nil
}
