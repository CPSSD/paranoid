package dnetclient

import (
	"time"
)

// Node struct containing the node information
type Node struct {
	IP   string
	Port string
}

// Nodes array
var Nodes []Node

// ThisNode has to be set before calling Join
var ThisNode Node

// DiscoveryAddr public string
var DiscoveryAddr string

var ResetInterval time.Duration
