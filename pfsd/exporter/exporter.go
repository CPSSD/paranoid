// exporter package creates an WebServer which exports Raft information
package exporter

import (
	"github.com/cpssd/paranoid/logger"
)

var (
	Log      *logger.ParanoidLogger
	server   *Server
	nodeList map[string]MessageNode
)

func init() {
	nodeList = make(map[string]MessageNode)
}

func Send(msg Message) {
	server.Send(msg)
}

func NewStdServer(port string) {
	server = NewServer(port)
}

func Listen() {
	server.Run()
}

// SetState is used to set the local exporter state
func SetState(nodes []MessageNode) {
	for i := 0; i < len(nodes); i++ {
		nodeList[nodes[i].Uuid] = nodes[i]
	}
}

// NodeChange sends a message to the client whenever there is a change in nodes
func NodeChange(node MessageNode) {
	msg := Message{
		Type: NodeChangeMessage,
		Data: MessageData{
			Node: node,
		},
	}

	// Check does the message exist to determine the response
	// TODO: 	Add implementation when a node is deleted. This is currently
	// 				not supported with our raft implementation
	if _, ok := nodeList[node.Uuid]; !ok {
		msg.Data.Action = "add"
	} else {
		msg.Data.Action = "update"
	}

	nodeList[node.Uuid] = node
	server.Send(msg)
}

// Send an event that happened to the client
func Event(msg Message) {
	server.Send(msg)
}

func toNodeArray(m map[string]MessageNode) []MessageNode {
	var nodes []MessageNode
	for _, n := range m {
		nodes = append(nodes, n)
	}
	return nodes
}
