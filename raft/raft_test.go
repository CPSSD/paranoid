// +build !integration

package raft

import (
	"fmt"
	"github.com/cpssd/paranoid/logger"
	pb "github.com/cpssd/paranoid/proto/raft"
	"log"
	"net"
	"os"
	"path"
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
	highestTerm := uint64(0)
	highestIndex := -1
	for i := 0; i < len(cluster); i++ {
		if isLeader(cluster[i]) {
			currentTerm := cluster[i].state.GetCurrentTerm()
			if currentTerm > highestTerm {
				highestTerm = currentTerm
				highestIndex = i
			}
			return cluster[i]
		}
	}
	if highestIndex > 0 {
		return cluster[highestIndex]
	}
	return nil
}

func closeListener(lis *net.Listener) {
	if lis != nil {
		(*lis).Close()
	}
}

func stopRaftServer(raftServer *RaftNetworkServer) {
	if raftServer.QuitChannelClosed == false {
		close(raftServer.Quit)
	}
}

func createPersistentFile(persistentFile string) string {
	dir, _ := path.Split(persistentFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0700)
		if err != nil {
			Log.Fatal("Error creating persistent file:", err)
		}
	}
	return persistentFile
}

func removePersistentFile(persistentFile string) {
	os.Remove(persistentFile)
}

func TestRaftElection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short testing mode")
	}

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

	node1PersistentPath := createPersistentFile(path.Join(os.TempDir(), "rafttest1", "node1"))
	defer removePersistentFile(node1PersistentPath)
	node1RaftServer, node1srv := startRaft(node1Lis, node1, node1PersistentPath, []Node{node2, node3})
	defer node1srv.Stop()
	defer stopRaftServer(node1RaftServer)

	node2PersistentPath := createPersistentFile(path.Join(os.TempDir(), "rafttest1", "node2"))
	defer removePersistentFile(node2PersistentPath)
	node2RaftServer, node2srv := startRaft(node2Lis, node2, node2PersistentPath, []Node{node1, node3})
	defer node2srv.Stop()
	defer stopRaftServer(node2RaftServer)

	node3PersistentPath := createPersistentFile(path.Join(os.TempDir(), "rafttest1", "node3"))
	defer removePersistentFile(node3PersistentPath)
	node3RaftServer, node3srv := startRaft(node3Lis, node3, node3PersistentPath, []Node{node1, node2})
	defer node3srv.Stop()
	defer stopRaftServer(node3RaftServer)

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
	if leader.state.nodeId == "node1" {
		node1srv.Stop()
	} else if leader.state.nodeId == "node2" {
		node2srv.Stop()
	} else {
		node3srv.Stop()
	}
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

func verifySpecialNumber(raftServer *RaftNetworkServer, x uint64, waitIntervals int) error {
	if raftServer.state.GetSpecialNumber() == x {
		return nil
	}
	for i := 0; i < waitIntervals; i++ {
		time.Sleep(500 * time.Millisecond)
		if raftServer.state.GetSpecialNumber() == x {
			return nil
		}
	}
	return fmt.Errorf(raftServer.state.nodeId, " special number", raftServer.state.GetSpecialNumber(), " is not equal to", x)
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

	node1PersistentPath := createPersistentFile(path.Join(os.TempDir(), "rafttest2", "node1"))
	defer removePersistentFile(node1PersistentPath)
	node1RaftServer, node1srv := startRaft(node1Lis, node1, node1PersistentPath, []Node{node2, node3})
	defer node1srv.Stop()
	defer stopRaftServer(node1RaftServer)

	node2PersistentPath := createPersistentFile(path.Join(os.TempDir(), "rafttest2", "node2"))
	defer removePersistentFile(node2PersistentPath)
	node2RaftServer, node2srv := startRaft(node2Lis, node2, node2PersistentPath, []Node{node1, node3})
	defer node2srv.Stop()
	defer stopRaftServer(node2RaftServer)

	node3PersistentPath := createPersistentFile(path.Join(os.TempDir(), "rafttest2", "node3"))
	defer removePersistentFile(node3PersistentPath)
	node3RaftServer, node3srv := startRaft(node3Lis, node3, node3PersistentPath, []Node{node1, node2})
	defer node3srv.Stop()
	defer stopRaftServer(node3RaftServer)

	err := node1RaftServer.RequestAddLogEntry(&pb.Entry{pb.Entry_StateMachineCommand, &pb.StateMachineCommand{10}, nil})
	cluster := []*RaftNetworkServer{node1RaftServer, node2RaftServer, node3RaftServer}
	leader := getLeader(cluster)

	if err != nil {
		Log.Info("most recent index :", node1RaftServer.state.log.GetMostRecentIndex())
		Log.Info("most recent leader index:", leader.state.log.GetMostRecentIndex())
		Log.Info("commit index:", leader.state.GetCommitIndex())
		Log.Info("leader commit:", leader.state.GetCommitIndex())
		t.Fatal("Failed to replicate entry,", err)
	}
	err = verifySpecialNumber(node1RaftServer, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	err = verifySpecialNumber(node2RaftServer, 10, 10)
	if err != nil {
		t.Fatal(err)
	}
	err = verifySpecialNumber(node3RaftServer, 10, 10)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRaftPersistentState(t *testing.T) {
	Log.Info("Testing persistent state")
	node1Lis, node1Port := startListener()
	node1 := setUpNode("node1", "localhost", node1Port, "_")
	defer closeListener(node1Lis)

	node1PersistentPath := createPersistentFile(path.Join(os.TempDir(), "rafttest2", "node1"))
	defer removePersistentFile(node1PersistentPath)
	node1RaftServer, node1srv := startRaft(node1Lis, node1, node1PersistentPath, []Node{})
	defer node1srv.Stop()
	defer stopRaftServer(node1RaftServer)

	err := node1RaftServer.RequestAddLogEntry(&pb.Entry{pb.Entry_StateMachineCommand, &pb.StateMachineCommand{10}, nil})
	if err != nil {
		t.Fatal("Test setup failed,", err)
	}

	cluster := []*RaftNetworkServer{node1RaftServer}

	count := 0
	var leader *RaftNetworkServer
	for {
		count++
		if count > 5 {
			t.Fatal("Test setup failed: Failed to select leader")
		}
		time.Sleep(1 * time.Second)
		leader = getLeader(cluster)
		if leader != nil {
			break
		}
	}

	close(node1RaftServer.Quit)
	node1srv.Stop()
	time.Sleep(1 * time.Second)

	currentTerm := node1RaftServer.state.GetCurrentTerm()
	Log.Info("Current Term:", currentTerm)
	lastApplied := node1RaftServer.state.GetLastApplied()
	Log.Info("Last applied:", lastApplied)
	votedFor := node1RaftServer.state.GetVotedFor()
	Log.Info("Voted For:", votedFor)

	node1RebootLis, _ := startListener()
	defer closeListener(node1RebootLis)

	node1RebootRaftServer, node1Rebootsrv := startRaft(node1RebootLis, node1, node1PersistentPath, []Node{})
	defer node1Rebootsrv.Stop()
	defer stopRaftServer(node1RebootRaftServer)

	if node1RebootRaftServer.state.GetCurrentTerm() != currentTerm {
		t.Fatal("Current term not restored after reboot. CurrentTerm:", node1RebootRaftServer.state.GetCurrentTerm())
	}
	if node1RebootRaftServer.state.GetLastApplied() != lastApplied {
		t.Fatal("Last applied not restored after reboot. Last applied:", node1RebootRaftServer.state.GetLastApplied())
	}
	if node1RebootRaftServer.state.GetVotedFor() != votedFor {
		t.Fatal("Voted for not restored after reboot. Last applied:", node1RebootRaftServer.state.GetVotedFor())
	}
}
