// +build !integration

package test

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	"github.com/cpssd/paranoid/raft"
	"github.com/cpssd/paranoid/raft/rafttestutil"
	"os"
	"path"
	"testing"
)

func TestSnapshoting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short testing mode")
	}
	t.Parallel()

	raft.Log.Info("Testing snapshoting")
	lis, node1Port := rafttestutil.StartListener()
	defer rafttestutil.CloseListener(lis)
	node := rafttestutil.SetUpNode("node", "localhost", node1Port, "_")
	raft.Log.Info("Listeners set up")

	raftDirectory := rafttestutil.CreateRaftDirectory(path.Join(os.TempDir(), "snapshottest", "node"))
	pfsDirectory := path.Join(os.TempDir(), "snapshottestpfs")
	err := os.RemoveAll(pfsDirectory)
	err = os.Mkdir(pfsDirectory, 0700)
	if err != nil {
		t.Fatal("Unable to make pfsdirectory:", err)
	}
	defer func() {
		err = os.RemoveAll(pfsDirectory)
		if err != nil {
			t.Fatal("Error removing pfsdirectory:", err)
		}
	}()

	code, err := commands.InitCommand(pfsDirectory)
	if code != returncodes.OK {
		t.Fatal("Unable to init pfsdirectroy:", err)
	}

	var raftServer *raft.RaftNetworkServer
	defer rafttestutil.RemoveRaftDirectory(raftDirectory, raftServer)
	raftServer, srv := raft.StartRaft(lis, node, pfsDirectory, raftDirectory, &raft.StartConfiguration{Peers: []raft.Node{}})
	defer srv.Stop()
	defer rafttestutil.StopRaftServer(raftServer)

	code, err = raftServer.RequestCreatCommand("test.txt", 0700)
	if code != returncodes.OK {
		t.Fatal("Error performing create command:", err)
	}

	code, err, _ = raftServer.RequestWriteCommand("test.txt", 0, 5, []byte("hello"))
	if code != returncodes.OK {
		t.Fatal("Error performing write command:", err)
	}

	//Take snapshot
	raftServer.Wait.Add(1)
	err = raftServer.CreateSnapshot(raftServer.State.Log.GetMostRecentIndex())
	if err != nil {
		t.Fatal("Error taking snapshot:", err)
	}

	code, err, _ = raftServer.RequestWriteCommand("test.txt", 0, 7, []byte("goodbye"))
	if code != returncodes.OK {
		t.Fatal("Error performing write command:", err)
	}

	//Revert to snapshot
	err = raftServer.RevertToSnapshot(path.Join(raftDirectory, raft.SnapshotDirectory, raft.CurrentSnapshotDirectory))
	if err != nil {
		t.Fatal("Error reverting to snapshot:", err)
	}

	code, err, data := commands.ReadCommand(pfsDirectory, "test.txt", -1, -1)
	if string(data) != "hello" {
		t.Fatal("Error reverting snapshot. Read does not match 'hello'. Actual:", string(data))
	}
}
