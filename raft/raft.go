package raft

import (
	"errors"
	"github.com/cpssd/paranoid/logger"
	pb "github.com/cpssd/paranoid/proto/raft"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	ELECTION_TIMEOUT      time.Duration = 3000 * time.Millisecond
	HEARTBEAT             time.Duration = 1000 * time.Millisecond
	REQUEST_VOTE_TIMEOUT  time.Duration = 5500 * time.Millisecond
	HEARTBEAT_TIMEOUT     time.Duration = 3000 * time.Millisecond
	SEND_ENTRY_TIMEOUT    time.Duration = 7500 * time.Millisecond
	ENTRY_APPLIED_TIMEOUT time.Duration = 20000 * time.Millisecond
)

var (
	Log *logger.ParanoidLogger
)

type RaftNetworkServer struct {
	State *RaftState
	Wait  sync.WaitGroup

	QuitChannelClosed    bool
	Quit                 chan bool
	ElectionTimeoutReset chan bool

	addEntryLock  sync.Mutex
	clientRequest *pb.Entry
}

func (s *RaftNetworkServer) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	if s.State.Configuration.InConfiguration(req.LeaderId) == false {
		if s.State.Configuration.InConfiguration(s.State.NodeId) == false {
			return &pb.AppendEntriesResponse{s.State.GetCurrentTerm(), false}, nil
		}
	}

	if req.Term < s.State.GetCurrentTerm() {
		return &pb.AppendEntriesResponse{s.State.GetCurrentTerm(), false}, nil
	}

	s.ElectionTimeoutReset <- true
	s.State.SetLeaderId(req.LeaderId)

	if req.Term > s.State.GetCurrentTerm() {
		s.State.SetCurrentTerm(req.Term)
		s.State.SetCurrentState(FOLLOWER)
	}

	if req.PrevLogIndex != 0 {
		preLogEntry := s.State.Log.GetLogEntry(req.PrevLogIndex)
		if preLogEntry == nil || preLogEntry.Term != req.PrevLogTerm {
			return &pb.AppendEntriesResponse{s.State.GetCurrentTerm(), false}, nil
		}
	}

	for i := uint64(0); i < uint64(len(req.Entries)); i++ {
		LogIndex := req.PrevLogIndex + 1 + i
		LogEntryAtIndex := s.State.Log.GetLogEntry(LogIndex)
		if LogEntryAtIndex != nil {
			if LogEntryAtIndex.Term != req.Term {
				s.State.Log.DiscardLogEntries(LogIndex)
				s.appendLogEntry(req.Entries[i])
			}
		} else {
			s.appendLogEntry(req.Entries[i])
		}
	}

	if req.LeaderCommit > s.State.GetCommitIndex() {
		lastLogIndex := s.State.Log.GetMostRecentIndex()
		if lastLogIndex < req.LeaderCommit {
			s.State.SetCommitIndex(lastLogIndex)
		} else {
			s.State.SetCommitIndex(req.LeaderCommit)
		}
		s.State.ApplyLogEntries()
	}

	return &pb.AppendEntriesResponse{s.State.GetCurrentTerm(), true}, nil
}

func (s *RaftNetworkServer) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	if s.State.Configuration.InConfiguration(req.CandidateId) == false {
		if s.State.Configuration.InConfiguration(s.State.NodeId) == false {
			return &pb.RequestVoteResponse{s.State.GetCurrentTerm(), false}, nil
		}
	}

	if req.Term < s.State.GetCurrentTerm() {
		return &pb.RequestVoteResponse{s.State.GetCurrentTerm(), false}, nil
	}

	if s.State.GetCurrentState() != CANDIDATE {
		return &pb.RequestVoteResponse{s.State.GetCurrentTerm(), false}, nil
	}

	if req.Term > s.State.GetCurrentTerm() {
		s.State.SetCurrentTerm(req.Term)
		s.State.SetCurrentState(FOLLOWER)
	}

	if req.LastLogTerm < s.State.Log.GetMostRecentTerm() {
		return &pb.RequestVoteResponse{s.State.GetCurrentTerm(), false}, nil
	} else if req.LastLogTerm == s.State.Log.GetMostRecentTerm() {
		if req.LastLogIndex < s.State.Log.GetMostRecentIndex() {
			return &pb.RequestVoteResponse{s.State.GetCurrentTerm(), false}, nil
		}
	}

	if s.State.GetVotedFor() == "" || s.State.GetVotedFor() == req.CandidateId {
		s.State.SetVotedFor(req.CandidateId)
		return &pb.RequestVoteResponse{s.State.GetCurrentTerm(), true}, nil
	}

	return &pb.RequestVoteResponse{s.State.GetCurrentTerm(), false}, nil
}

func (s *RaftNetworkServer) getLeader() *Node {
	leaderId := s.State.GetLeaderId()
	if leaderId != "" {
		if s.State.Configuration.InConfiguration(leaderId) {
			node, err := s.State.Configuration.GetNode(leaderId)
			if err == nil {
				return &node
			} else {
				return nil
			}
		}
	}
	return nil
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

func (s *RaftNetworkServer) appendLogEntry(entry *pb.Entry) {
	if entry.Type == pb.Entry_ConfigurationChange {
		config := entry.GetConfig()
		if config == nil {
			Log.Fatal("Incorrect entry information. No configuration present")
		}
		if config.Type == pb.Configuration_CurrentConfiguration {
			s.State.Configuration.UpdateCurrentConfiguration(protoNodesToNodes(config.Nodes), s.State.Log.GetMostRecentIndex())
		} else {
			s.State.Configuration.NewFutureConfiguration(protoNodesToNodes(config.Nodes), s.State.Log.GetMostRecentIndex())
		}
	}
	s.State.Log.AppendEntry(entry, s.State.GetCurrentTerm())
}

func (s *RaftNetworkServer) addLogEntryLeader(entry *pb.Entry) error {
	if entry.Type == pb.Entry_ConfigurationChange {
		config := entry.GetConfig()
		if config != nil {
			if config.Type == pb.Configuration_FutureConfiguration {
				if s.State.Configuration.GetFutureConfigurationActive() {
					return errors.New("Can not change confirugation while another configuration change is underway")
				}
			}
		} else {
			return errors.New("Incorrect entry information. No configuration present")
		}
	}
	s.appendLogEntry(entry)
	s.State.calculateNewCommitIndex()
	s.State.SendAppendEntries <- true
	return nil
}

func (s *RaftNetworkServer) ClientToLeaderRequest(ctx context.Context, req *pb.EntryRequest) (*pb.EmptyMessage, error) {
	if s.State.Configuration.InConfiguration(req.SenderId) == false {
		return &pb.EmptyMessage{}, errors.New("Node is not in the configuration")
	}

	if s.State.GetCurrentState() != LEADER {
		return &pb.EmptyMessage{}, errors.New("Node is not the current leader")
	}
	err := s.addLogEntryLeader(req.Entry)
	return &pb.EmptyMessage{}, err
}

func (s *RaftNetworkServer) sendLeaderLogEntry(entry *pb.Entry) error {
	leaderNode := s.getLeader()
	if leaderNode == nil {
		return errors.New("Unable to find leader")
	}

	conn, err := Dial(leaderNode, SEND_ENTRY_TIMEOUT)
	defer conn.Close()
	if err == nil {
		client := pb.NewRaftNetworkClient(conn)
		_, err := client.ClientToLeaderRequest(context.Background(), &pb.EntryRequest{s.State.NodeId, entry})
		return err
	}
	return err
}

//A request to add a Log entry from a client. If the node is not the leader, it must forward the request to the leader.
//Only return once the request has been commited to the State machine
func (s *RaftNetworkServer) RequestAddLogEntry(entry *pb.Entry) error {
	s.addEntryLock.Lock()
	defer s.addEntryLock.Unlock()
	currentState := s.State.GetCurrentState()

	s.State.SetWaitingForApplied(true)
	defer s.State.SetWaitingForApplied(false)

	//Add entry to leaders Log
	if currentState == LEADER {
		err := s.addLogEntryLeader(entry)
		if err != nil {
			return err
		}
	} else if currentState == FOLLOWER {
		if s.State.GetLeaderId() != "" {
			err := s.sendLeaderLogEntry(entry)
			if err != nil {
				return err
			}
		} else {
			select {
			case <-time.After(20 * time.Second):
				return errors.New("Could not find a leader")
			case <-s.State.LeaderElected:
				if s.State.GetCurrentState() == LEADER {
					err := s.addLogEntryLeader(entry)
					if err != nil {
						return err
					}
				} else {
					err := s.sendLeaderLogEntry(entry)
					if err != nil {
						return err
					}
				}
			}
		}
	} else {
		count := 0
		for {
			count++
			if count > 40 {
				return errors.New("Could not find a leader")
			}
			time.Sleep(500 * time.Millisecond)
			currentState = s.State.GetCurrentState()
			if currentState != CANDIDATE {
				break
			}
		}
		if currentState == LEADER {
			err := s.addLogEntryLeader(entry)
			if err != nil {
				return err
			}
		} else {
			err := s.sendLeaderLogEntry(entry)
			if err != nil {
				return err
			}
		}
	}

	//Wait for the Log entry to be applied
	timer := time.NewTimer(ENTRY_APPLIED_TIMEOUT)
	for {
		select {
		case <-timer.C:
			return errors.New("Waited too long to commit Log entry")
		case entryIndex := <-s.State.EntryApplied:
			LogEntry := s.State.Log.GetLogEntry(entryIndex)
			if LogEntry.Entry.Uuid == entry.Uuid {
				return nil
			}
		}
	}
	return nil
}

func generateNewUUID() string {
	uuidBytes, err := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	if err != nil {
		Log.Fatal("Error generating new UUID:", err)
	}
	return strings.TrimSpace(string(uuidBytes))
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

func (s *RaftNetworkServer) RequestChangeConfiguration(nodes []Node) error {
	Log.Info("Configuration change requested")
	entry := &pb.Entry{
		Type:    pb.Entry_ConfigurationChange,
		Uuid:    generateNewUUID(),
		Command: nil,
		Config: &pb.Configuration{
			Type:  pb.Configuration_FutureConfiguration,
			Nodes: convertNodesToProto(nodes),
		},
	}
	return s.RequestAddLogEntry(entry)
}

func getRandomElectionTimeout() time.Duration {
	rand.Seed(time.Now().UnixNano())
	return ELECTION_TIMEOUT + time.Duration(rand.Intn(int(ELECTION_TIMEOUT)))
}

func (s *RaftNetworkServer) electionTimeOut() {
	timer := time.NewTimer(getRandomElectionTimeout())
	defer s.Wait.Done()
	defer timer.Stop()
	for {
		select {
		case _, ok := <-s.Quit:
			if !ok {
				s.QuitChannelClosed = true
				Log.Info("Exiting election timeout loop")
				return
			}
		case <-s.ElectionTimeoutReset:
			timer.Reset(getRandomElectionTimeout())
		case <-timer.C:
			Log.Info("Starting new election")
			s.State.SetCurrentTerm(s.State.GetCurrentTerm() + 1)
			s.State.SetCurrentState(CANDIDATE)
			timer.Reset(getRandomElectionTimeout())
		}
	}
}

func Dial(node *Node, timeoutMiliseconds time.Duration) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTimeout(timeoutMiliseconds))
	//TODO: tls support
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(node.String(), opts...)
	return conn, err
}

func (s *RaftNetworkServer) requestPeerVote(node *Node, term uint64, voteChannel chan *voteResponse) {
	defer s.Wait.Done()
	for {
		if term != s.State.GetCurrentTerm() || s.State.GetCurrentState() != CANDIDATE {
			voteChannel <- nil
			return
		}
		Log.Info("Dialing ", node)
		conn, err := Dial(node, REQUEST_VOTE_TIMEOUT)
		defer conn.Close()
		if err == nil {
			client := pb.NewRaftNetworkClient(conn)
			response, err := client.RequestVote(context.Background(), &pb.RequestVoteRequest{s.State.GetCurrentTerm(),
				s.State.NodeId,
				s.State.Log.GetMostRecentIndex(),
				s.State.Log.GetMostRecentTerm()})
			Log.Info("Got response from", node)
			if err == nil {
				voteChannel <- &voteResponse{response, node.NodeID}
				return
			}
		}
	}
}

type voteResponse struct {
	response *pb.RequestVoteResponse
	NodeId   string
}

func (s *RaftNetworkServer) runElection() {
	defer s.Wait.Done()
	term := s.State.GetCurrentTerm()
	var votesGranted []string

	if s.State.GetVotedFor() == "" && s.State.Configuration.InConfiguration(s.State.NodeId) {
		s.State.SetVotedFor(s.State.NodeId)
		votesGranted = append(votesGranted, s.State.NodeId)
	}

	if s.State.Configuration.HasMajority(votesGranted) {
		Log.Info("Node elected leader with", len(votesGranted), " votes")
		s.State.SetCurrentState(LEADER)
		return
	}

	Log.Info("Sending RequestVote RPCs to peers")
	voteChannel := make(chan *voteResponse)
	peers := s.State.Configuration.GetPeersList()
	for i := 0; i < len(peers); i++ {
		s.Wait.Add(1)
		go s.requestPeerVote(&peers[i], term, voteChannel)
	}

	votesReturned := 0
	votesRequested := len(peers)
	for {
		select {
		case _, ok := <-s.Quit:
			if !ok {
				s.QuitChannelClosed = true
				Log.Info("Exiting election loop")
				return
			}
		case vote := <-voteChannel:
			if term != s.State.GetCurrentTerm() || s.State.GetCurrentState() != CANDIDATE {
				return
			}
			votesReturned++
			if vote != nil {
				if vote.response.Term > s.State.GetCurrentTerm() {
					Log.Info("Stopping election, higher term encountered.")
					s.State.SetCurrentTerm(vote.response.Term)
					s.State.SetCurrentState(FOLLOWER)
					return
				} else {
					if vote.response.VoteGranted == true {
						votesGranted = append(votesGranted, vote.NodeId)
						Log.Info("Vote granted. Current votes :", len(votesGranted))
						if s.State.Configuration.HasMajority(votesGranted) {
							Log.Info("Node elected leader with", len(votesGranted), " votes")
							s.State.SetCurrentState(LEADER)
							return
						}
					}
				}
			}
			if votesReturned == votesRequested {
				return
			}
		}
	}
}

func (s *RaftNetworkServer) manageElections() {
	defer s.Wait.Done()
	for {
		select {
		case _, ok := <-s.Quit:
			if !ok {
				s.QuitChannelClosed = true
				Log.Info("Exiting elction managment loop")
				return
			}
		case <-s.State.StartElection:
			s.Wait.Add(1)
			go s.runElection()
		}
	}
}

func (s *RaftNetworkServer) sendHeartBeat(node *Node) {
	defer s.Wait.Done()
	nextIndex := s.State.Configuration.GetNextIndex(node.NodeID)
	conn, err := Dial(node, HEARTBEAT_TIMEOUT)
	defer conn.Close()
	if err == nil {
		client := pb.NewRaftNetworkClient(conn)
		if s.State.Log.GetMostRecentIndex() >= nextIndex {
			prevLogEntry := s.State.Log.GetLogEntry(nextIndex - 1)
			prevLogTerm := uint64(0)
			if prevLogEntry != nil {
				prevLogTerm = prevLogEntry.Term
			}

			response, err := client.AppendEntries(context.Background(), &pb.AppendEntriesRequest{
				Term:         s.State.GetCurrentTerm(),
				LeaderId:     s.State.NodeId,
				PrevLogIndex: nextIndex - 1,
				PrevLogTerm:  prevLogTerm,
				Entries:      []*pb.Entry{&s.State.Log.GetLogEntry(nextIndex).Entry},
				LeaderCommit: s.State.GetCommitIndex(),
			})
			if err == nil {
				if response.Term > s.State.GetCurrentTerm() {
					s.State.StopLeading <- true
				} else if response.Success == false {
					if s.State.GetCurrentState() == LEADER {
						s.State.Configuration.SetNextIndex(node.NodeID, nextIndex-1)
					}
				} else if response.Success {
					if s.State.GetCurrentState() == LEADER {
						s.State.Configuration.SetNextIndex(node.NodeID, nextIndex+1)
						s.State.Configuration.SetMatchIndex(node.NodeID, nextIndex)
						s.State.calculateNewCommitIndex()
					}
				}
			}
		} else {
			response, err := client.AppendEntries(context.Background(), &pb.AppendEntriesRequest{
				Term:         s.State.GetCurrentTerm(),
				LeaderId:     s.State.NodeId,
				PrevLogIndex: s.State.Log.GetMostRecentIndex(),
				PrevLogTerm:  s.State.Log.GetMostRecentTerm(),
				Entries:      []*pb.Entry{},
				LeaderCommit: s.State.GetCommitIndex(),
			})
			if err == nil {
				if response.Term > s.State.GetCurrentTerm() {
					s.State.StopLeading <- true
				}
			}
		}
	}
}

func (s *RaftNetworkServer) manageLeading() {
	defer s.Wait.Done()
	for {
		select {
		case <-s.State.StopLeading:
		case <-s.State.SendAppendEntries:
		case _, ok := <-s.Quit:
			if !ok {
				s.QuitChannelClosed = true
				s.State.SetCurrentState(INACTIVE)
				Log.Info("Exiting leading managment loop")
				return
			}
		case <-s.State.StartLeading:
			Log.Info("Started leading for term ", s.State.GetCurrentTerm())
			s.State.Configuration.ResetNodeIndexs(s.State.Log.GetMostRecentIndex())
			peers := s.State.Configuration.GetPeersList()
			for i := 0; i < len(peers); i++ {
				s.Wait.Add(1)
				go s.sendHeartBeat(&peers[i])
			}
			timer := time.NewTimer(HEARTBEAT)
			for {
				select {
				case _, ok := <-s.Quit:
					if !ok {
						s.QuitChannelClosed = true
						Log.Info("Exiting heartbeat loop")
						return
					}
				case <-s.State.StopLeading:
					Log.Info("Stopped leading")
					return
				case <-s.State.SendAppendEntries:
					timer.Reset(HEARTBEAT)
					s.ElectionTimeoutReset <- true
					peers = s.State.Configuration.GetPeersList()
					for i := 0; i < len(peers); i++ {
						s.Wait.Add(1)
						go s.sendHeartBeat(&peers[i])
					}
				case <-timer.C:
					timer.Reset(HEARTBEAT)
					s.ElectionTimeoutReset <- true
					peers = s.State.Configuration.GetPeersList()
					for i := 0; i < len(peers); i++ {
						s.Wait.Add(1)
						go s.sendHeartBeat(&peers[i])
					}
				}
			}
		}
	}
}

func (s *RaftNetworkServer) manageConfigurationChanges() {
	defer s.Wait.Done()
	for {
		select {
		case _, ok := <-s.Quit:
			if !ok {
				s.QuitChannelClosed = true
				Log.Info("Exiting configuration managment loop")
				return
			}
		case config := <-s.State.ConfigurationApplied:
			if config.Type == pb.Configuration_CurrentConfiguration {
				inConfig := false
				for i := 0; i < len(config.Nodes); i++ {
					if config.Nodes[i].NodeId == s.State.NodeId {
						inConfig = true
						break
					}
				}
				if inConfig == false {
					Log.Info("Node not included in current configuration", s.State.NodeId)
					s.State.SetCurrentState(FOLLOWER)
				}
			} else {
				if s.State.GetCurrentState() == LEADER {
					newConfig := &pb.Entry{
						Type:    pb.Entry_ConfigurationChange,
						Uuid:    generateNewUUID(),
						Command: nil,
						Config: &pb.Configuration{
							Type:  pb.Configuration_CurrentConfiguration,
							Nodes: config.Nodes,
						},
					}
					s.addLogEntryLeader(newConfig)
				}
			}
		}
	}
}

func newRaftNetworkServer(nodeDetails Node, persistentStateFile string, peers []Node) *RaftNetworkServer {
	raftServer := &RaftNetworkServer{State: newRaftState(nodeDetails, persistentStateFile, peers)}
	raftServer.ElectionTimeoutReset = make(chan bool, 100)
	raftServer.Quit = make(chan bool)
	raftServer.QuitChannelClosed = false
	raftServer.Wait.Add(4)
	go raftServer.electionTimeOut()
	go raftServer.manageElections()
	go raftServer.manageLeading()
	go raftServer.manageConfigurationChanges()
	return raftServer
}

func StartRaft(lis *net.Listener, nodeDetails Node, persistentStateFile string, peers []Node) (*RaftNetworkServer, *grpc.Server) {
	var opts []grpc.ServerOption
	srv := grpc.NewServer(opts...)
	raftServer := newRaftNetworkServer(nodeDetails, persistentStateFile, peers)
	pb.RegisterRaftNetworkServer(srv, raftServer)
	raftServer.Wait.Add(1)
	go func() {
		Log.Info("RaftNetworkServer started")
		err := srv.Serve(*lis)
		if err != nil {
			Log.Error("Error running RaftNetworkServer", err)
		}
	}()
	return raftServer, srv
}
