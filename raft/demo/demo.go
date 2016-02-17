package main

import (
	"flag"
	"github.com/cpssd/paranoid/logger"
	pb "github.com/cpssd/paranoid/proto/raft"
	"github.com/cpssd/paranoid/raft"
	"github.com/cpssd/paranoid/raft/rafttestutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

const (
	DEMO_DURATION         time.Duration = 50 * time.Second
	PRINT_TIME            time.Duration = 5 * time.Second
	RANDOM_NUMBER_GEN_MIN time.Duration = 3000 * time.Millisecond
	RANDOM_NUMBER_GEN_MAX time.Duration = 10000 * time.Millisecond
)

var (
	waitGroup sync.WaitGroup
	demo      = flag.Int("demo", 0, "Which demo to run (1-3). 0 for all demos")
)

func getRandomInterval() time.Duration {
	interval := int(RANDOM_NUMBER_GEN_MAX) - int(RANDOM_NUMBER_GEN_MIN)
	randx := rand.Intn(interval)
	return RANDOM_NUMBER_GEN_MIN + time.Duration(randx)
}

func manageNode(raftServer *raft.RaftNetworkServer) {
	defer waitGroup.Done()
	testTimer := time.NewTimer(DEMO_DURATION)
	randomNumTimer := time.NewTimer(getRandomInterval())
	for {
		select {
		case <-testTimer.C:
			return
		case <-randomNumTimer.C:
			if raftServer.QuitChannelClosed {
				return
			}
			randomNumber := rand.Intn(1000)
			log.Println(raftServer.State.NodeId, "requesting that", randomNumber, "be added to the log")
			err := raftServer.RequestAddLogEntry(&pb.Entry{
				pb.Entry_StateMachineCommand,
				rafttestutil.GenerateNewUUID(),
				&pb.StateMachineCommand{uint64(randomNumber)},
				nil,
			})
			if err == nil {
				log.Println(raftServer.State.NodeId, "successfullly added", randomNumber, "to the log")
			} else {
				log.Println(raftServer.State.NodeId, "could not add", randomNumber, "to the log:", err)
			}
			randomNumTimer.Reset(getRandomInterval())
		}
	}
}

func printLogs(cluster []*raft.RaftNetworkServer) {
	defer waitGroup.Done()
	testTimer := time.NewTimer(DEMO_DURATION)
	printTimer := time.NewTimer(PRINT_TIME)
	for {
		select {
		case <-testTimer.C:
			return
		case <-printTimer.C:
			printTimer.Reset(PRINT_TIME)
			log.Println("Printing node logs:")
			for i := 0; i < len(cluster); i++ {
				logsString := ""
				for j := uint64(1); j <= cluster[i].State.Log.GetMostRecentIndex(); j++ {
					logsString = logsString + " " + strconv.Itoa(int(cluster[i].State.Log.GetLogEntry(j).Entry.GetCommand().Number))
				}
				log.Println(cluster[i].State.NodeId, "Logs:", logsString)
			}
		}
	}
}

func performDemoOne(node1RaftServer, node2RaftServer, node3RaftServer *raft.RaftNetworkServer) {
	log.Println("Running basic raft demo")
	waitGroup.Add(4)
	go manageNode(node1RaftServer)
	go manageNode(node2RaftServer)
	go manageNode(node3RaftServer)
	go printLogs([]*raft.RaftNetworkServer{node1RaftServer, node2RaftServer, node3RaftServer})

	waitGroup.Wait()
}

func performDemoTwo(node1srv *grpc.Server, node1RaftServer, node2RaftServer, node3RaftServer *raft.RaftNetworkServer) {
	log.Println("Running raft node crash demo")
	waitGroup.Add(4)
	go manageNode(node1RaftServer)
	go manageNode(node2RaftServer)
	go manageNode(node3RaftServer)
	go printLogs([]*raft.RaftNetworkServer{node1RaftServer, node2RaftServer, node3RaftServer})

	//Node1 will crash after 20 seconds
	time.Sleep(20 * time.Second)
	log.Println("Crashing node1")
	node1srv.Stop()
	rafttestutil.StopRaftServer(node1RaftServer)

	waitGroup.Wait()
}

func setupDemo(demoNum int) {
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

	if demoNum == 1 {
		performDemoOne(node1RaftServer, node2RaftServer, node3RaftServer)
	} else if demoNum == 2 {
		performDemoTwo(node1srv, node1RaftServer, node2RaftServer, node3RaftServer)
	}
}

func createDemoDirectory() {
	removeDemoDirectory()
	err := os.Mkdir(path.Join(os.TempDir(), "raftdemo"), 0777)
	if err != nil {
		log.Fatalln("Error creating demo directory:", err)
	}
}

func removeDemoDirectory() {
	os.RemoveAll(path.Join(os.TempDir(), "raftdemo"))
}

func createFileWriter(logPath string, component string) (io.Writer, error) {
	return os.OpenFile(path.Join(logPath, component+".log"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	createDemoDirectory()
	defer removeDemoDirectory()
	raft.Log = logger.New("raftdemo", "raftdemo", path.Join(os.TempDir(), "raftdemo"))
	err := raft.Log.SetOutput(logger.LOGFILE)
	if err != nil {
		log.Println("Could not set file logging:", err)
	}

	writer, err := createFileWriter(path.Join(os.TempDir(), "raftdemo"), "grpclog")
	if err != nil {
		log.Println("Could not redirect grpc output")
	} else {
		grpclog.SetLogger(log.New(writer, "", log.LstdFlags))
	}

	flag.Parse()

	demo := *demo
	if demo == 0 {
		setupDemo(1)
		setupDemo(2)
		setupDemo(3)
	} else {
		if demo > 3 || demo < 1 {
			log.Fatal("Bad demo number specified")
		}
		setupDemo(demo)
	}
}
