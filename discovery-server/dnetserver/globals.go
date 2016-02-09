// Package dnetserver implements the DiscoveryNetwork gRPC server.
// globals.go contains data used by each gRPC handler in dnetserver.
package dnetserver

import (
	"github.com/cpssd/paranoid/logger"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"time"
)

var Log *logger.ParanoidLogger

// DiscoveryServer struct
type DiscoveryServer struct{}

// Node struct to hold the node data
type Node struct {
	Pool       string
	ExpiryTime time.Time
	Data       pb.Node
}

// Nodes array
var Nodes []Node

// RenewInterval global containing the time after which the nodes will be marked as inactive
var RenewInterval time.Duration
