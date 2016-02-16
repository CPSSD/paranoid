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
	ELECTION_TIMEOUT      time.Duration = 3000
	HEARTBEAT             time.Duration = 1000
	REQUEST_VOTE_TIMEOUT  time.Duration = 5500
	HEARTBEAT_TIMEOUT     time.Duration = 3000
	SEND_ENTRY_TIMEOUT    time.Duration = 7500
	ENTRY_APPLIED_TIMEOUT time.Duration = 20000
)

var (
	Log *logger.ParanoidLogger
)

type RaftNetworkServer struct {
	state *RaftState
	Wait  sync.WaitGroup

	QuitChannelClosed    bool
	Quit                 chan bool
	ElectionTimeoutReset chan bool

	addEntryLock  sync.Mutex
	clientRequest *pb.Entry
}

func (s *RaftNetworkServer) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	if req.Term < s.state.GetCurrentTerm() {
		return &pb.AppendEntriesResponse{s.state.GetCurrentTerm(), false}, nil
	}

	s.ElectionTimeoutReset <- true
	s.state.SetLeaderId(req.LeaderId)

	if req.Term > s.state.GetCurrentTerm() {
		s.state.SetCurrentTerm(req.Term)
		s.state.SetCurrentState(FOLLOWER)
	}

	if req.PrevLogIndex != 0 {
		preLogEntry := s.state.log.GetLogEntry(req.PrevLogIndex)
		if preLogEntry == nil || preLogEntry.Term != req.PrevLogTerm {
			return &pb.AppendEntriesResponse{s.state.GetCurrentTerm(), false}, nil
		}
	}

	for i := uint64(0); i < uint64(len(req.Entries)); i++ {
		logIndex := req.PrevLogIndex + 1 + i
		logEntryAtIndex := s.state.log.GetLogEntry(logIndex)
		if logEntryAtIndex != nil {
			if logEntryAtIndex.Term != req.Term {
				s.state.log.DiscardLogEntries(logIndex)
				s.appendLogEntry(req.Entries[i])
			}
		} else {
			s.appendLogEntry(req.Entries[i])
		}
	}

	if req.LeaderCommit > s.state.GetCommitIndex() {
		lastLogIndex := s.state.log.GetMostRecentIndex()
		if lastLogIndex < req.LeaderCommit {
			s.state.SetCommitIndex(lastLogIndex)
		} else {
			s.state.SetCommitIndex(req.LeaderCommit)
		}
		s.state.ApplyLogEntries()
	}

	return &pb.AppendEntriesResponse{s.state.GetCurrentTerm(), true}, nil
}

func (s *RaftNetworkServer) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	if req.Term < s.state.GetCurrentTerm() {
		return &pb.RequestVoteResponse{s.state.GetCurrentTerm(), false}, nil
	}

	if s.state.GetCurrentState() != CANDIDATE {
		return &pb.RequestVoteResponse{s.state.GetCurrentTerm(), false}, nil
	}

	if req.Term > s.state.GetCurrentTerm() {
		s.state.SetCurrentTerm(req.Term)
		s.state.SetCurrentState(FOLLOWER)
	}

	if req.LastLogTerm < s.state.log.GetMostRecentTerm() {
		return &pb.RequestVoteResponse{s.state.GetCurrentTerm(), false}, nil
	} else if req.LastLogTerm == s.state.log.GetMostRecentTerm() {
		if req.LastLogIndex < s.state.log.GetMostRecentIndex() {
			return &pb.RequestVoteResponse{s.state.GetCurrentTerm(), false}, nil
		}
	}

	if s.state.GetVotedFor() == "" || s.state.GetVotedFor() == req.CandidateId {
		s.state.SetVotedFor(req.CandidateId)
		return &pb.RequestVoteResponse{s.state.GetCurrentTerm(), true}, nil
	}

	return &pb.RequestVoteResponse{s.state.GetCurrentTerm(), false}, nil
}

func (s *RaftNetworkServer) getLeader() *Node {
	leaderId := s.state.GetLeaderId()
	if leaderId != "" {
		node := s.state.Configuration.GetNode(leaderId)
		return &node
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
			s.state.Configuration.UpdateCurrentConfiguration(protoNodesToNodes(config.Nodes), s.state.log.GetMostRecentIndex())
		} else {
			s.state.Configuration.NewFutureConfiguration(protoNodesToNodes(config.Nodes), s.state.log.GetMostRecentIndex())
		}
	}
	s.state.log.AppendEntry(entry, s.state.GetCurrentTerm())
}

func (s *RaftNetworkServer) addLogEntryLeader(entry *pb.Entry) error {
	if entry.Type == pb.Entry_ConfigurationChange {
		config := entry.GetConfig()
		if config != nil {
			if config.Type == pb.Configuration_FutureConfiguration {
				if s.state.Configuration.GetFutureConfigurationActive() {
					return errors.New("Can not change confirugation while another configuration change is underway")
				}
			}
		} else {
			return errors.New("Incorrect entry information. No configuration present")
		}
	}
	s.appendLogEntry(entry)
	s.state.calculateNewCommitIndex()
	s.state.SendAppendEntries <- true
	return nil
}

func (s *RaftNetworkServer) ClientToLeaderRequest(ctx context.Context, req *pb.EntryRequest) (*pb.EmptyMessage, error) {
	if s.state.GetCurrentState() != LEADER {
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
		_, err := client.ClientToLeaderRequest(context.Background(), &pb.EntryRequest{entry})
		return err
	}
	return err
}

//A request to add a log entry from a client. If the node is not the leader, it must forward the request to the leader.
//Only return once the request has been commited to the state machine
func (s *RaftNetworkServer) RequestAddLogEntry(entry *pb.Entry) error {
	s.addEntryLock.Lock()
	defer s.addEntryLock.Unlock()
	currentState := s.state.GetCurrentState()

	s.state.SetWaitingForApplied(true)
	defer s.state.SetWaitingForApplied(false)

	//Add entry to leaders log
	if currentState == LEADER {
		err := s.addLogEntryLeader(entry)
		if err != nil {
			return err
		}
	} else if currentState == FOLLOWER {
		if s.state.GetLeaderId() != "" {
			err := s.sendLeaderLogEntry(entry)
			if err != nil {
				return err
			}
		} else {
			select {
			case <-time.After(20 * time.Second):
				return errors.New("Could not find a leader")
			case <-s.state.LeaderElected:
				if s.state.GetCurrentState() == LEADER {
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
			if count > 20 {
				return errors.New("Could not find a leader")
			}
			time.Sleep(500 * time.Millisecond)
			currentState = s.state.GetCurrentState()
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

	//Wait for the log entry to be applied
	timer := time.NewTimer(ENTRY_APPLIED_TIMEOUT * time.Millisecond)
	for {
		select {
		case <-timer.C:
			return errors.New("Waited too long to commit log entry")
		case entryIndex := <-s.state.EntryApplied:
			logEntry := s.state.log.GetLogEntry(entryIndex)
			if logEntry.Entry.Uuid == entry.Uuid {
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
	timer := time.NewTimer(getRandomElectionTimeout() * time.Millisecond)
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
			timer.Reset(getRandomElectionTimeout() * time.Millisecond)
		case <-timer.C:
			Log.Info("Starting new election")
			s.state.SetCurrentTerm(s.state.GetCurrentTerm() + 1)
			s.state.SetCurrentState(CANDIDATE)
			timer.Reset(getRandomElectionTimeout() * time.Millisecond)
		}
	}
}

func Dial(node *Node, timeoutMiliseconds time.Duration) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTimeout(timeoutMiliseconds*time.Millisecond))
	//TODO: tls support
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(node.String(), opts...)
	return conn, err
}

func (s *RaftNetworkServer) requestPeerVote(node *Node, term uint64, voteChannel chan *voteResponse) {
	defer s.Wait.Done()
	for {
		if term != s.state.GetCurrentTerm() || s.state.GetCurrentState() != CANDIDATE {
			voteChannel <- nil
			return
		}
		Log.Info("Dialing ", node)
		conn, err := Dial(node, REQUEST_VOTE_TIMEOUT)
		defer conn.Close()
		if err == nil {
			client := pb.NewRaftNetworkClient(conn)
			response, err := client.RequestVote(context.Background(), &pb.RequestVoteRequest{s.state.GetCurrentTerm(),
				s.state.nodeId,
				s.state.log.GetMostRecentIndex(),
				s.state.log.GetMostRecentTerm()})
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
	nodeId   string
}

func (s *RaftNetworkServer) runElection() {
	defer s.Wait.Done()
	term := s.state.GetCurrentTerm()
	var votesGranted []string

	if s.state.GetVotedFor() == "" && s.state.Configuration.InConfiguration(s.state.nodeId) {
		s.state.SetVotedFor(s.state.nodeId)
		votesGranted = append(votesGranted, s.state.nodeId)
	}

	if s.state.Configuration.HasMajority(votesGranted) {
		Log.Info("Node elected leader with", len(votesGranted), " votes")
		s.state.SetCurrentState(LEADER)
		return
	}

	Log.Info("Sending RequestVote RPCs to peers")
	voteChannel := make(chan *voteResponse)
	peers := s.state.Configuration.GetPeersList()
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
			if term != s.state.GetCurrentTerm() || s.state.GetCurrentState() != CANDIDATE {
				return
			}
			votesReturned++
			if vote != nil {
				if vote.response.Term > s.state.GetCurrentTerm() {
					Log.Info("Stopping election, higher term encountered.")
					s.state.SetCurrentTerm(vote.response.Term)
					s.state.SetCurrentState(FOLLOWER)
					return
				} else {
					if vote.response.VoteGranted == true {
						votesGranted = append(votesGranted, vote.nodeId)
						Log.Info("Vote granted. Current votes :", len(votesGranted))
						if s.state.Configuration.HasMajority(votesGranted) {
							Log.Info("Node elected leader with", len(votesGranted), " votes")
							s.state.SetCurrentState(LEADER)
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
		case <-s.state.StartElection:
			s.Wait.Add(1)
			go s.runElection()
		}
	}
}

func (s *RaftNetworkServer) sendHeartBeat(node *Node) {
	defer s.Wait.Done()
	nextIndex := s.state.Configuration.GetNextIndex(node.NodeID)
	conn, err := Dial(node, HEARTBEAT_TIMEOUT)
	defer conn.Close()
	if err == nil {
		client := pb.NewRaftNetworkClient(conn)
		if s.state.log.GetMostRecentIndex() >= nextIndex {
			prevLogEntry := s.state.log.GetLogEntry(nextIndex - 1)
			prevLogTerm := uint64(0)
			if prevLogEntry != nil {
				prevLogTerm = prevLogEntry.Term
			}

			response, err := client.AppendEntries(context.Background(), &pb.AppendEntriesRequest{
				Term:         s.state.GetCurrentTerm(),
				LeaderId:     s.state.nodeId,
				PrevLogIndex: nextIndex - 1,
				PrevLogTerm:  prevLogTerm,
				Entries:      []*pb.Entry{&s.state.log.GetLogEntry(nextIndex).Entry},
				LeaderCommit: s.state.GetCommitIndex(),
			})
			if err == nil {
				if response.Term > s.state.GetCurrentTerm() {
					s.state.StopLeading <- true
				} else if response.Success == false {
					if s.state.GetCurrentState() == LEADER {
						s.state.Configuration.SetNextIndex(node.NodeID, nextIndex-1)
					}
				} else if response.Success {
					if s.state.GetCurrentState() == LEADER {
						s.state.Configuration.SetNextIndex(node.NodeID, nextIndex+1)
						s.state.Configuration.SetMatchIndex(node.NodeID, nextIndex)
						s.state.calculateNewCommitIndex()
					}
				}
			}
		} else {
			response, err := client.AppendEntries(context.Background(), &pb.AppendEntriesRequest{
				Term:         s.state.GetCurrentTerm(),
				LeaderId:     s.state.nodeId,
				PrevLogIndex: s.state.log.GetMostRecentIndex(),
				PrevLogTerm:  s.state.log.GetMostRecentTerm(),
				Entries:      []*pb.Entry{},
				LeaderCommit: s.state.GetCommitIndex(),
			})
			if err == nil {
				if response.Term > s.state.GetCurrentTerm() {
					s.state.StopLeading <- true
				}
			}
		}
	}
}

func (s *RaftNetworkServer) manageLeading() {
	defer s.Wait.Done()
	for {
		select {
		case <-s.state.StopLeading:
		case <-s.state.SendAppendEntries:
		case _, ok := <-s.Quit:
			if !ok {
				s.QuitChannelClosed = true
				Log.Info("Exiting leading managment loop")
				return
			}
		case <-s.state.StartLeading:
			Log.Info("Started leading for term ", s.state.GetCurrentTerm())
			s.state.Configuration.ResetNodeIndexs(s.state.log.GetMostRecentIndex())
			peers := s.state.Configuration.GetPeersList()
			for i := 0; i < len(peers); i++ {
				s.Wait.Add(1)
				go s.sendHeartBeat(&peers[i])
			}
			timer := time.NewTimer(HEARTBEAT * time.Millisecond)
			for {
				select {
				case _, ok := <-s.Quit:
					if !ok {
						s.QuitChannelClosed = true
						Log.Info("Exiting heartbeat loop")
						return
					}
				case <-s.state.StopLeading:
					break
				case <-s.state.SendAppendEntries:
					timer.Reset(HEARTBEAT * time.Millisecond)
					s.ElectionTimeoutReset <- true
					peers = s.state.Configuration.GetPeersList()
					for i := 0; i < len(peers); i++ {
						s.Wait.Add(1)
						go s.sendHeartBeat(&peers[i])
					}
				case <-timer.C:
					timer.Reset(HEARTBEAT * time.Millisecond)
					s.ElectionTimeoutReset <- true
					peers = s.state.Configuration.GetPeersList()
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
		case config := <-s.state.ConfigurationApplied:
			if config.Type == pb.Configuration_CurrentConfiguration {
				inConfig := false
				for i := 0; i < len(config.Nodes); i++ {
					if config.Nodes[i].NodeId == s.state.nodeId {
						inConfig = true
						break
					}
				}
				if inConfig == false {
					s.state.SetCurrentState(FOLLOWER)
				}
			} else {
				if s.state.GetCurrentState() == LEADER {
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
	raftServer := &RaftNetworkServer{state: newRaftState(nodeDetails, persistentStateFile, peers)}
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

func startRaft(lis *net.Listener, nodeDetails Node, persistentStateFile string, peers []Node) (*RaftNetworkServer, *grpc.Server) {
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
