// +build !integration
package raft

import (
	"github.com/cpssd/paranoid/logger"
	"log"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	Log = logger.New("rafttest", "rafttest", os.DevNull)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func startListener() (*net.Listener, string) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Failed to start listening : %v.\n", err)
	}
	splits := strings.Split(lis.Addr().String(), ":")
	port := splits[len(splits)-1]
	return &lis, port
}

func setUpNode(name, ip, port, commonName string) Node {
	return Node{
		NodeID:     name,
		IP:         ip,
		Port:       port,
		CommonName: commonName,
	}
}

func isLeader(server *RaftNetworkServer) bool {
	return server.state.GetCurrentState() == LEADER
}

func TestRaftElection(t *testing.T) {
	node1Lis, node1Port := startListener()
	node1 := setUpNode("node1", "localhost", node1Port, "_")
	node2Lis, node2Port := startListener()
	node2 := setUpNode("node2", "localhost", node2Port, "_")
	node3Lis, node3Port := startListener()
	node3 := setUpNode("node3", "localhost", node3Port, "_")
	Log.Info("Listeners set up")

	node1RaftServer, node1srv := startRaft(node1Lis, "node1", []Node{node2, node3})
	node2RaftServer, node2srv := startRaft(node2Lis, "node2", []Node{node1, node3})
	node3RaftServer, node3srv := startRaft(node3Lis, "node3", []Node{node1, node2})
	defer node1srv.Stop()
	defer node2srv.Stop()
	defer node3srv.Stop()

	Log.Info("Searching for leader")
	count := 0
	for {
		count++
		if count > 5 {
			t.Error("Failed to select leader")
			break
		}
		time.Sleep(5 * time.Second)
		if isLeader(node1RaftServer) {
			t.Log("Node1 selected as leader for term", node1RaftServer.state.GetCurrentTerm())
			break
		}
		if isLeader(node2RaftServer) {
			t.Log("Node2 selected as leader for term", node2RaftServer.state.GetCurrentTerm())
			break
		}
		if isLeader(node3RaftServer) {
			t.Log("Node3 selected as leader for term", node3RaftServer.state.GetCurrentTerm())
			break
		}
	}
}
