// +build !integration
package raft

import (
	"github.com/cpssd/paranoid/logger"
	pb "github.com/cpssd/paranoid/proto/raft"
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

func closeListener(lis *net.Listener) {
	if lis != nil {
		(*lis).Close()
	}
}

func TestRaftElection(t *testing.T) {
	Log.Info("Testing leader eleciton")
	node1Lis, node1Port := startListener()
	defer closeListener(node1Lis)
	node1 := setUpNode("node1", "localhost", node1Port, "_")
	node2Lis, node2Port := startListener()
	defer closeListener(node2Lis)
	node2 := setUpNode("node2", "localhost", node2Port, "_")
	node3Lis, node3Port := startListener()
	defer closeListener(node3Lis)
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

	//Shutdown current leader, make sure an election is triggered and another leader is found
	close(leader.Quit)
	time.Sleep(5 * time.Second)

	for {
		count++
		if count > 5 {
			t.Fatal("Failed to select leader after original leader is shut down")
		}
		time.Sleep(5 * time.Second)
		newleader := getLeader(cluster)
		if leader == newleader {
			t.Fatal("Old leader failed to shut down")
		}
		if leader != nil {
			t.Log(newleader.state.nodeId, "selected as leader for term", leader.state.GetCurrentTerm())
			break
		}
	}
}

func TestRaftLogReplication(t *testing.T) {
	Log.Info("Testing log replication")
	node1Lis, node1Port := startListener()
	defer closeListener(node1Lis)
	node1 := setUpNode("node1", "localhost", node1Port, "_")
	node2Lis, node2Port := startListener()
	defer closeListener(node2Lis)
	node2 := setUpNode("node2", "localhost", node2Port, "_")
	node3Lis, node3Port := startListener()
	defer closeListener(node3Lis)
	node3 := setUpNode("node3", "localhost", node3Port, "_")
	Log.Info("Listeners set up")

	node1RaftServer, node1srv := startRaft(node1Lis, "node1", []Node{node2, node3})
	_, node2srv := startRaft(node2Lis, "node2", []Node{node1, node3})
	_, node3srv := startRaft(node3Lis, "node3", []Node{node1, node2})
	defer node1srv.Stop()
	defer node2srv.Stop()
	defer node3srv.Stop()

	err := node1RaftServer.RequestAddLogEntry(&pb.Entry{10})
	if err == nil {
		t.Error("Failed to replicate entry,", err)
	}
}
