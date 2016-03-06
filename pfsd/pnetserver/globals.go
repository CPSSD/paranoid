// Package pnetserver implements the ParanoidNetwork gRPC server.
// globals.go contains data used by each gRPC handler in pnetserver.
package pnetserver

import (
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/raft"
)

type ParanoidServer struct{}

var Log *logger.ParanoidLogger

var RaftNetworkServer *raft.RaftNetworkServer
