package raft

import (
	"github.com/cpssd/paranoid/pfsd/exporter"
	pb "github.com/cpssd/paranoid/proto/raft"
	"io/ioutil"
	"strings"
)

func protoDetailedNodeToExportNode(nodes []*pb.LeaderData_Data_DetailedNode) []exporter.MessageNode {
	res := make([]exporter.MessageNode, len(nodes))
	for i := 0; i < len(nodes); i++ {
		res[i] = exporter.MessageNode{
			CommonName: nodes[i].CommonName,
			Addr:       nodes[i].Addr,
			Uuid:       nodes[i].Uuid,
			State:      nodes[i].State,
		}
	}
	return res
}

func protoNodesToNodes(protoNodes []*pb.Node) []Node {
	nodes := make([]Node, len(protoNodes))
	for i := 0; i < len(protoNodes); i++ {
		nodes[i] = Node{
			IP:         protoNodes[i].Ip,
			Port:       protoNodes[i].Port,
			CommonName: protoNodes[i].CommonName,
			NodeID:     protoNodes[i].NodeId,
		}
	}
	return nodes
}

func convertNodesToProto(nodes []Node) []*pb.Node {
	protoNodes := make([]*pb.Node, len(nodes))
	for i := 0; i < len(nodes); i++ {
		protoNodes[i] = &pb.Node{
			Ip:         nodes[i].IP,
			Port:       nodes[i].Port,
			CommonName: nodes[i].CommonName,
			NodeId:     nodes[i].NodeID,
		}
	}
	return protoNodes
}

func generateNewUUID() string {
	uuidBytes, err := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	if err != nil {
		Log.Fatal("Error generating new UUID:", err)
	}
	return strings.TrimSpace(string(uuidBytes))
}
