package rafttestutil

import (
	"github.com/cpssd/paranoid/raft"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func GenerateNewUUID() string {
	uuidBytes, err := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	if err != nil {
		log.Fatalln("Error generating new UUID:", err)
	}
	return strings.TrimSpace(string(uuidBytes))
}

func StartListener() (*net.Listener, string) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Failed to start listening : %v.\n", err)
	}
	splits := strings.Split(lis.Addr().String(), ":")
	port := splits[len(splits)-1]
	return &lis, port
}

func SetUpNode(name, ip, port, commonName string) raft.Node {
	return raft.Node{
		NodeID:     name,
		IP:         ip,
		Port:       port,
		CommonName: commonName,
	}
}

func CloseListener(lis *net.Listener) {
	if lis != nil {
		(*lis).Close()
		file, _ := (*lis).(*net.TCPListener).File()
		file.Close()
	}
}

func StopRaftServer(raftServer *raft.RaftNetworkServer) {
	if raftServer.QuitChannelClosed == false {
		close(raftServer.Quit)
	}
}

func CreateRaftDirectory(raftDirectory string) string {
	os.RemoveAll(raftDirectory)
	err := os.MkdirAll(raftDirectory, 0700)
	if err != nil {
		log.Fatal("Error creating raft directory:", err)
	}
	return raftDirectory
}

func RemoveRaftDirectory(raftDirectory string) {
	//Need to sleep, as otherwise this can cause a directory that the tests are using to be removed before
	//the raft servers have shut down. Causing the tests to fail.
	time.Sleep(time.Second)
	os.RemoveAll(raftDirectory)
}

func IsLeader(server *raft.RaftNetworkServer) bool {
	return server.State.GetCurrentState() == raft.LEADER
}

func GetLeader(cluster []*raft.RaftNetworkServer) *raft.RaftNetworkServer {
	highestTerm := uint64(0)
	highestIndex := -1
	for i := 0; i < len(cluster); i++ {
		if IsLeader(cluster[i]) {
			currentTerm := cluster[i].State.GetCurrentTerm()
			if currentTerm > highestTerm {
				highestTerm = currentTerm
				highestIndex = i
			}
		}
	}
	if highestIndex >= 0 {
		return cluster[highestIndex]
	}
	return nil
}

func GetLeaderTimeout(cluster []*raft.RaftNetworkServer, timeoutSeconds int) *raft.RaftNetworkServer {
	var leader *raft.RaftNetworkServer
	leader = GetLeader(cluster)
	if leader != nil {
		return leader
	}
	count := 0
	for {
		count++
		if count > timeoutSeconds {
			break
		}
		time.Sleep(1 * time.Second)
		leader = GetLeader(cluster)
		if leader != nil {
			break
		}
	}
	return leader
}
