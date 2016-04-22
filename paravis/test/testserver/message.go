package main

// Type is the message type
type Type string

const (
	// TypeState will contain the current state of the network
	TypeState Type = "state"
	// TypeNodeChange will be send whenever a node changes in the cluster
	TypeNodeChange = "nodechange"
	// TypeEvent is send when something happens, like a write
	TypeEvent = "event"
)

// NodeState is the state that the node could be
type NodeState string

const (
	// Follower is the basic type
	Follower NodeState = "follower"
	// Current is who you are
	Current = "current"
	// Leader is the leader of the cluster
	Leader = "leader"
	// Inactive nodes are still in cluser but.. inactive
	Inactive = "inactive"
)

func (ns NodeState) String() string {
	return string(ns)
}

// Node contains the node information
type Node struct {
	CommonName string `json:"commonName"`
	UUID       string `json:"uuid"`
	Addr       string `json:"addr"`
	State      string `json:"state"`
}

// Event contains all the event data
type Event struct {
	Source  string `json:"source"`
	Target  string `json:"target"`
	Details string `json:"details"`
}

// Data contains the data of the message
type Data struct {
	Nodes  []Node `json:"nodes,omitempty"`
	Action string `json:"action,omitempty"`
	Node   Node   `json:"node,omitempty"`
	Event  Event  `json:"event,omitempty"`
}

// Message is the message being send to the client
type Message struct {
	Type Type `json:"type"`
	Data Data `json:"data"`
}
