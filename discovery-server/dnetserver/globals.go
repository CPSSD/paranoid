// Package dnetserver implements the DiscoveryNetwork gRPC server.
// globals.go contains data used by each gRPC handler in dnetserver.
package dnetserver

import (
	"github.com/cpssd/paranoid/logger"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"sync"
	"time"
)

var Log *logger.ParanoidLogger

// DiscoveryServer struct
type DiscoveryServer struct{}

// PoolInfo struct to hold the pool data
type PoolInfo struct {
	PasswordSalt []byte `json:"passwordsalt"`
	PasswordHash []byte `json:"passwordhash"`
	Nodes        map[string]*pb.Node
}

type Pool struct {
	PoolLock sync.Mutex
	Info     PoolInfo
}

// PoolLock is locked when acessing the pool map
var PoolLock sync.RWMutex

// Pools map
var Pools map[string]*Pool

// RenewInterval global containing the time after which the nodes will be marked as inactive
var RenewInterval time.Duration

// StateDirectoryPath is the path to the directory in which the discovery server stores its state
var StateDirectoryPath string

// TempDirecotryPath is the path to the directory where temporary state files are stored
var TempDirectoryPath string

func checkPoolPassword(pool, password string, node *pb.Node) error {
	if Pools[pool] != nil {
		if password == "" {
			if len(Pools[pool].Info.PasswordHash) != 0 {
				Log.Errorf("Join: node %s attempted join password protected pool without a giving a password", node.Uuid)
				returnError := grpc.Errorf(codes.Internal,
					"pool %s is password protected",
					pool,
				)
				return returnError
			}
		} else {
			err := bcrypt.CompareHashAndPassword(Pools[pool].Info.PasswordHash, append(Pools[pool].Info.PasswordSalt, []byte(password)...))
			if err != nil {
				Log.Errorf("Join: node %s attempted join password protected pool with incorrect password: %s",
					node.Uuid,
					err,
				)
				returnError := grpc.Errorf(codes.Internal,
					"given password incorrect: %s",
					err,
				)
				return returnError
			}
		}
	}
	return nil
}
