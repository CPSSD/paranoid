// +build !integration

package test

import (
	"fmt"
	"github.com/cpssd/paranoid/logger"
	pb "github.com/cpssd/paranoid/proto/raft"
	"github.com/cpssd/paranoid/raft"
	"github.com/cpssd/paranoid/raft/rafttestutil"
	"os"
	"path"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	raft.Log = logger.New("rafttest", "rafttest", os.DevNull)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestRaftElection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short testing mode")
	}

	raft.Log.Info("Testing leader eleciton")
	node1Lis, node1Port := rafttestutil.StartListener()
	defer rafttestutil.CloseListener(node1Lis)
	node1 := rafttestutil.SetUpNode("node1", "localhost", node1Port, "_")
	node2Lis, node2Port := rafttestutil.StartListener()
	defer rafttestutil.CloseListener(node2Lis)
	node2 := rafttestutil.SetUpNode("node2", "localhost", node2Port, "_")
	node3Lis, node3Port := rafttestutil.StartListener()
	defer rafttestutil.CloseListener(node3Lis)
	node3 := rafttestutil.SetUpNode("node3", "localhost", node3Port, "_")
	raft.Log.Info("Listeners set up")

	node1PersistentPath := rafttestutil.CreatePersistentFile(path.Join(os.TempDir(), "rafttest1", "node1"))
	defer rafttestutil.RemovePersistentFile(node1PersistentPath)
	node1RaftServer, node1srv := raft.StartRaft(node1Lis, node1, node1PersistentPath, []raft.Node{node2, node3})
	defer node1srv.Stop()
	defer rafttestutil.StopRaftServer(node1RaftServer)

	node2PersistentPath := rafttestutil.CreatePersistentFile(path.Join(os.TempDir(), "rafttest1", "node2"))
	defer rafttestutil.RemovePersistentFile(node2PersistentPath)
	node2RaftServer, node2srv := raft.StartRaft(node2Lis, node2, node2PersistentPath, []raft.Node{node1, node3})
	defer node2srv.Stop()
	defer rafttestutil.StopRaftServer(node2RaftServer)

	node3PersistentPath := rafttestutil.CreatePersistentFile(path.Join(os.TempDir(), "rafttest1", "node3"))
	defer rafttestutil.RemovePersistentFile(node3PersistentPath)
	node3RaftServer, node3srv := raft.StartRaft(node3Lis, node3, node3PersistentPath, []raft.Node{node1, node2})
	defer node3srv.Stop()
	defer rafttestutil.StopRaftServer(node3RaftServer)

	cluster := []*raft.RaftNetworkServer{node1RaftServer, node2RaftServer, node3RaftServer}

	raft.Log.Info("Searching for leader")
	leader := rafttestutil.GetLeaderTimeout(cluster, 25)
	if leader != nil {
		t.Log(leader.State.NodeId, "selected as leader for term", leader.State.GetCurrentTerm())
	} else {
		t.Fatal("Failed to select leader")
	}

	//Shutdown current leader, make sure an election is triggered and another leader is found
	close(leader.Quit)
	if leader.State.NodeId == "node1" {
		node1srv.Stop()
	} else if leader.State.NodeId == "node2" {
		node2srv.Stop()
	} else {
		node3srv.Stop()
	}
	time.Sleep(5 * time.Second)

	count := 0
	for {
		count++
		if count > 5 {
			t.Fatal("Failed to select leader after original leader is shut down")
		}
		time.Sleep(5 * time.Second)
		newLeader := rafttestutil.GetLeader(cluster)
		if newLeader != nil && leader != newLeader {
			t.Log(newLeader.State.NodeId, "selected as leader for term", newLeader.State.GetCurrentTerm())
			break
		}
	}
}

func verifySpecialNumber(raftServer *raft.RaftNetworkServer, x uint64, waitIntervals int) error {
	if raftServer.State.GetSpecialNumber() == x {
		return nil
	}
	for i := 0; i < waitIntervals; i++ {
		time.Sleep(500 * time.Millisecond)
		if raftServer.State.GetSpecialNumber() == x {
			return nil
		}
	}
	return fmt.Errorf(raftServer.State.NodeId, " special number", raftServer.State.GetSpecialNumber(), " is not equal to", x)
}

func TestRaftLogReplication(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	raft.Log.Info("Testing log replication")
	node1Lis, node1Port := rafttestutil.StartListener()
	defer rafttestutil.CloseListener(node1Lis)
	node1 := rafttestutil.SetUpNode("node1", "localhost", node1Port, "_")
	node2Lis, node2Port := rafttestutil.StartListener()
	defer rafttestutil.CloseListener(node2Lis)
	node2 := rafttestutil.SetUpNode("node2", "localhost", node2Port, "_")
	node3Lis, node3Port := rafttestutil.StartListener()
	defer rafttestutil.CloseListener(node3Lis)
	node3 := rafttestutil.SetUpNode("node3", "localhost", node3Port, "_")
	raft.Log.Info("Listeners set up")

	node1PersistentPath := rafttestutil.CreatePersistentFile(path.Join(os.TempDir(), "rafttest2", "node1"))
	defer rafttestutil.RemovePersistentFile(node1PersistentPath)
	node1RaftServer, node1srv := raft.StartRaft(node1Lis, node1, node1PersistentPath, []raft.Node{node2, node3})
	defer node1srv.Stop()
	defer rafttestutil.StopRaftServer(node1RaftServer)

	node2PersistentPath := rafttestutil.CreatePersistentFile(path.Join(os.TempDir(), "rafttest2", "node2"))
	defer rafttestutil.RemovePersistentFile(node2PersistentPath)
	node2RaftServer, node2srv := raft.StartRaft(node2Lis, node2, node2PersistentPath, []raft.Node{node1, node3})
	defer node2srv.Stop()
	defer rafttestutil.StopRaftServer(node2RaftServer)

	node3PersistentPath := rafttestutil.CreatePersistentFile(path.Join(os.TempDir(), "rafttest2", "node3"))
	defer rafttestutil.RemovePersistentFile(node3PersistentPath)
	node3RaftServer, node3srv := raft.StartRaft(node3Lis, node3, node3PersistentPath, []raft.Node{node1, node2})
	defer node3srv.Stop()
	defer rafttestutil.StopRaftServer(node3RaftServer)

	err := node1RaftServer.RequestAddLogEntry(&pb.Entry{
		Type: pb.Entry_Demo,
		Uuid: rafttestutil.GenerateNewUUID(),
		Demo: &pb.DemoCommand{10},
	})
	cluster := []*raft.RaftNetworkServer{node1RaftServer, node2RaftServer, node3RaftServer}
	leader := rafttestutil.GetLeader(cluster)

	if err != nil {
		raft.Log.Info("most recent index :", node1RaftServer.State.Log.GetMostRecentIndex())
		raft.Log.Info("most recent leader index:", leader.State.Log.GetMostRecentIndex())
		raft.Log.Info("commit index:", leader.State.GetCommitIndex())
		raft.Log.Info("leader commit:", leader.State.GetCommitIndex())
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
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	raft.Log.Info("Testing persistent State")
	node1Lis, node1Port := rafttestutil.StartListener()
	node1 := rafttestutil.SetUpNode("node1", "localhost", node1Port, "_")
	defer rafttestutil.CloseListener(node1Lis)

	node1PersistentPath := rafttestutil.CreatePersistentFile(path.Join(os.TempDir(), "rafttest2", "node1"))
	defer rafttestutil.RemovePersistentFile(node1PersistentPath)
	node1RaftServer, node1srv := raft.StartRaft(node1Lis, node1, node1PersistentPath, []raft.Node{})
	defer node1srv.Stop()
	defer rafttestutil.StopRaftServer(node1RaftServer)

	err := node1RaftServer.RequestAddLogEntry(&pb.Entry{
		Type: pb.Entry_Demo,
		Uuid: rafttestutil.GenerateNewUUID(),
		Demo: &pb.DemoCommand{10},
	})
	if err != nil {
		t.Fatal("Test setup failed,", err)
	}

	cluster := []*raft.RaftNetworkServer{node1RaftServer}

	leader := rafttestutil.GetLeaderTimeout(cluster, 1)
	if leader == nil {
		t.Fatal("Test setup failed: Failed to select leader")
	}

	close(node1RaftServer.Quit)
	node1srv.Stop()
	time.Sleep(1 * time.Second)

	currentTerm := node1RaftServer.State.GetCurrentTerm()
	raft.Log.Info("Current Term:", currentTerm)
	lastApplied := node1RaftServer.State.GetLastApplied()
	raft.Log.Info("Last applied:", lastApplied)
	votedFor := node1RaftServer.State.GetVotedFor()
	raft.Log.Info("Voted For:", votedFor)

	node1RebootLis, _ := rafttestutil.StartListener()
	defer rafttestutil.CloseListener(node1RebootLis)

	node1RebootRaftServer, node1Rebootsrv := raft.StartRaft(node1RebootLis, node1, node1PersistentPath, []raft.Node{})
	defer node1Rebootsrv.Stop()
	defer rafttestutil.StopRaftServer(node1RebootRaftServer)

	if node1RebootRaftServer.State.GetCurrentTerm() != currentTerm {
		t.Fatal("Current term not restored after reboot. CurrentTerm:", node1RebootRaftServer.State.GetCurrentTerm())
	}
	if node1RebootRaftServer.State.GetLastApplied() != lastApplied {
		t.Fatal("Last applied not restored after reboot. Last applied:", node1RebootRaftServer.State.GetLastApplied())
	}
	if node1RebootRaftServer.State.GetVotedFor() != votedFor {
		t.Fatal("Voted for not restored after reboot. Last applied:", node1RebootRaftServer.State.GetVotedFor())
	}
}

func TestRaftConfigurationChange(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	raft.Log.Info("Testing joinging and leaving cluster")

	node1Lis, node1Port := rafttestutil.StartListener()
	defer rafttestutil.CloseListener(node1Lis)
	node1 := rafttestutil.SetUpNode("node1", "localhost", node1Port, "_")

	node2Lis, node2Port := rafttestutil.StartListener()
	defer rafttestutil.CloseListener(node2Lis)
	node2 := rafttestutil.SetUpNode("node2", "localhost", node2Port, "_")

	node3Lis, node3Port := rafttestutil.StartListener()
	defer rafttestutil.CloseListener(node3Lis)
	node3 := rafttestutil.SetUpNode("node3", "localhost", node3Port, "_")

	node4Lis, node4Port := rafttestutil.StartListener()
	defer rafttestutil.CloseListener(node4Lis)
	node4 := rafttestutil.SetUpNode("node4", "localhost", node4Port, "_")

	raft.Log.Info("Listeners set up")

	node1PersistentPath := rafttestutil.CreatePersistentFile(path.Join(os.TempDir(), "rafttest3", "node1"))
	defer rafttestutil.RemovePersistentFile(node1PersistentPath)
	node1RaftServer, node1srv := raft.StartRaft(node1Lis, node1, node1PersistentPath, []raft.Node{node2, node3})
	defer node1srv.Stop()
	defer rafttestutil.StopRaftServer(node1RaftServer)

	node2PersistentPath := rafttestutil.CreatePersistentFile(path.Join(os.TempDir(), "rafttest3", "node2"))
	defer rafttestutil.RemovePersistentFile(node2PersistentPath)
	node2RaftServer, node2srv := raft.StartRaft(node2Lis, node2, node2PersistentPath, []raft.Node{node1, node3})
	defer node2srv.Stop()
	defer rafttestutil.StopRaftServer(node2RaftServer)

	node3PersistentPath := rafttestutil.CreatePersistentFile(path.Join(os.TempDir(), "rafttest3", "node3"))
	defer rafttestutil.RemovePersistentFile(node3PersistentPath)
	node3RaftServer, node3srv := raft.StartRaft(node3Lis, node3, node3PersistentPath, []raft.Node{node1, node2})
	defer node3srv.Stop()
	defer rafttestutil.StopRaftServer(node3RaftServer)

	node4PersistentPath := rafttestutil.CreatePersistentFile(path.Join(os.TempDir(), "rafttest3", "node4"))
	defer rafttestutil.RemovePersistentFile(node4PersistentPath)
	node4RaftServer, node4srv := raft.StartRaft(node4Lis, node4, node4PersistentPath, []raft.Node{})
	defer node4srv.Stop()
	defer rafttestutil.StopRaftServer(node4RaftServer)

	err := node1RaftServer.RequestAddLogEntry(&pb.Entry{
		Type: pb.Entry_Demo,
		Uuid: rafttestutil.GenerateNewUUID(),
		Demo: &pb.DemoCommand{10},
	})
	if err != nil {
		t.Fatal("Test setup failed:", err)
	}

	err = node2RaftServer.RequestChangeConfiguration([]raft.Node{node1, node2, node3, node4})
	if err != nil {
		t.Fatal("Unabled to change configuration:", err)
	}

	err = verifySpecialNumber(node4RaftServer, 10, 10)
	if err != nil {
		t.Fatal(err)
	}

	cluster := []*raft.RaftNetworkServer{node1RaftServer, node2RaftServer, node3RaftServer, node4RaftServer}

	leader := rafttestutil.GetLeaderTimeout(cluster, 15)
	if leader == nil {
		t.Fatal("Unable to find a leader")
	}

	raft.Log.Info(leader.State.NodeId, " is to leave configuration")

	newNodes := []raft.Node{node1, node2, node3, node4}
	for i := 0; i < len(newNodes); i++ {
		if newNodes[i].NodeID == leader.State.NodeId {
			newNodes = append(newNodes[:i], newNodes[i+1:]...)
			raft.Log.Info("Removing leader from new set of nodes")
			break
		}
	}

	err = node1RaftServer.RequestChangeConfiguration(newNodes)
	if err != nil {
		t.Fatal("Unable to change configuration:", err)
	}

	count := 0
	var newLeader *raft.RaftNetworkServer
	for {
		count++
		if count > 10 {
			t.Fatal("Old leader did not stepdown in time")
		}
		time.Sleep(2 * time.Second)
		newLeader = rafttestutil.GetLeader(cluster)
		if newLeader != nil && newLeader.State.NodeId != leader.State.NodeId {
			break
		}
	}

	time.Sleep(raft.HEARTBEAT_TIMEOUT * 2)

	if node1RaftServer.State.NodeId != leader.State.NodeId {
		err := node1RaftServer.RequestAddLogEntry(&pb.Entry{
			Type: pb.Entry_Demo,
			Uuid: rafttestutil.GenerateNewUUID(),
			Demo: &pb.DemoCommand{1337},
		})
		if err != nil {
			t.Fatal("Unable to commit new entry:", err)
		}
	} else {
		err := node2RaftServer.RequestAddLogEntry(&pb.Entry{
			Type: pb.Entry_StateMachineCommand,
			Uuid: rafttestutil.GenerateNewUUID(),
			Demo: &pb.DemoCommand{1337},
		})
		if err != nil {
			t.Fatal("Unable to commit new entry:", err)
		}
	}
}
