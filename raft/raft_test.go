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

func getLeader(cluster []*RaftNetworkServer) *RaftNetworkServer {
	for i := 0; i < len(cluster); i++ {
		if isLeader(cluster[i]) {
			return cluster[i]
		}
	}
	return nil
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

	cluster := []*RaftNetworkServer{node1RaftServer, node2RaftServer, node3RaftServer}

	Log.Info("Searching for leader")
	count := 0
	var leader *RaftNetworkServer
	for {
		count++
		if count > 5 {
			t.Fatal("Failed to select leader")
		}
		time.Sleep(5 * time.Second)
		leader = getLeader(cluster)
		if leader != nil {
			t.Log(leader.state.nodeId, "selected as leader for term", leader.state.GetCurrentTerm())
			break
		}
	}
}
