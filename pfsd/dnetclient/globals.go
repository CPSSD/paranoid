package dnetclient

import (
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
)

// Nodes array
var Nodes []*pb.Node

var thisNode pb.Node

var DiscoveryAddr string

var resetInterval int64
