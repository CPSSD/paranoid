package demo

import (
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/raft"
	"github.com/cpssd/paranoid/raft/rafttestutil"
	"log"
	"os"
	"path"
)

func setupDemo() {
	node1Lis, node1Port := rafttestutil.StartListener()
	defer rafttestutil.CloseListener(node1Lis)
	node1 := rafttestutil.SetUpNode("node1", "localhost", node1Port, "_")
	node2Lis, node2Port := rafttestutil.StartListener()
	defer rafttestutil.CloseListener(node2Lis)
	node2 := rafttestutil.SetUpNode("node2", "localhost", node2Port, "_")
	node3Lis, node3Port := rafttestutil.StartListener()
	defer rafttestutil.CloseListener(node3Lis)
	node3 := rafttestutil.SetUpNode("node3", "localhost", node3Port, "_")

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
}

func createDemoDirectory() {
	err := os.Mkdir(path.Join(os.TempDir(), "raftdemo"), 0777)
	if err != nil {
		log.Fatalln("Error creating demo directory:", err)
	}
}

func removeDemoDirectory() {
	err := os.RemoveAll(path.Join(os.TempDir(), "raftdemo"))
	if err != nil {
		log.Println("Could not delete demo directory:", err)
	}
}

func main() {
	createDemoDirectory()
	defer removeDemoDirectory()
	raft.Log = logger.New("raftdemo", "raftdemo", path.Join(os.TempDir(), "raftdemo", "logs"))
	raft.Log.SetOutput(logger.LOGFILE)

	setupDemo()
}
