package raft

import (
	"errors"
	"github.com/cpssd/paranoid/logger"
	pb "github.com/cpssd/paranoid/proto/raft"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	ELECTION_TIMEOUT     time.Duration = 3000
	HEARTBEAT            time.Duration = 1000
	REQUEST_VOTE_TIMEOUT time.Duration = 5500
	HEARTBEAT_TIMEOUT    time.Duration = 3000
	SEND_ENTRY_TIMEOUT   time.Duration = 7500
)

var (
	Log *logger.ParanoidLogger
)

type RaftNetworkServer struct {
	state *RaftState
	Wait  sync.WaitGroup

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
				s.state.log.AppendEntry(req.Entries[logIndex], req.Term)
			}
		} else {
			s.state.log.AppendEntry(req.Entries[i], req.Term)
		}
	}

	if req.LeaderCommit > s.state.GetCommitIndex() {
		lastLogIndex := s.state.log.GetMostRecentIndex()
		if lastLogIndex < req.LeaderCommit {
			s.state.SetCommitIndex(lastLogIndex)
		} else {
			s.state.SetCommitIndex(req.LeaderCommit)
		}
	}

	return &pb.AppendEntriesResponse{s.state.GetCurrentTerm(), true}, nil
}

func (s *RaftNetworkServer) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	if req.Term < s.state.GetCurrentTerm() {
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
	for i := 0; i < len(s.state.peers); i++ {
		if s.state.peers[i].NodeID == leaderId {
			return &s.state.peers[i]
		}
	}
	return nil
}

func (s *RaftNetworkServer) addLogEntry(entry *pb.Entry) {
	s.state.log.AppendEntry(entry, s.state.GetCurrentTerm())
	s.state.SendAppendEntries <- true
}

func (s *RaftNetworkServer) ClientToLeaderRequest(ctx context.Context, req *pb.EntryRequest) (*pb.EmptyMessage, error) {
	if s.state.GetCurrentState() != LEADER {
		return &pb.EmptyMessage{}, errors.New("Node is not the current leader")
	}
	s.addLogEntry(req.Entry)
	return &pb.EmptyMessage{}, nil
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
		s.addLogEntry(entry)
	} else if currentState == FOLLOWER {
		if s.state.GetLeaderId() != "" {
			err := s.sendLeaderLogEntry(entry)
			if err != nil {
				return err
			}
		} else {
			select {
			case <-time.After(10 * time.Second):
				return errors.New("Could not find a leader")
			case <-s.state.LeaderElected:
				err := s.sendLeaderLogEntry(entry)
				if err != nil {
					return err
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
			s.addLogEntry(entry)
		} else {
			err := s.sendLeaderLogEntry(entry)
			if err != nil {
				return err
			}
		}
	}

	//Wait for the log entry to be commited
	timer := time.NewTimer(20 * time.Second)
	for {
		select {
		case <-timer.C:
			return errors.New("Waited too long to commit log entry")
		case entryIndex := <-s.state.EntryApplied:
			logEntry := s.state.log.GetLogEntry(entryIndex)
			if logEntry.Entry == *entry {
				Log.Info("Entry sucessfully commited to log: ", logEntry.Entry)
				return nil
			}
		}
	}
	return nil
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

func (s *RaftNetworkServer) requestPeerVote(node *Node, term uint64, voteChannel chan *pb.RequestVoteResponse) {
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
				voteChannel <- response
				return
			}
		}
	}
}

func (s *RaftNetworkServer) getRequiredVotes() int {
	nodecount := 1 + len(s.state.peers)
	if nodecount%2 == 0 {
		return nodecount/2 + 1
	}
	return nodecount / 2
}

func (s *RaftNetworkServer) runElection() {
	defer s.Wait.Done()
	term := s.state.GetCurrentTerm()
	votesGranted := 0
	votesRequired := s.getRequiredVotes()

	if s.state.GetVotedFor() == "" {
		s.state.SetVotedFor(s.state.nodeId)
		votesGranted++
	}

	Log.Info("Sending RequestVote RPCs to peers")
	voteChannel := make(chan *pb.RequestVoteResponse)
	for i := 0; i < len(s.state.peers); i++ {
		s.Wait.Add(1)
		go s.requestPeerVote(&s.state.peers[i], term, voteChannel)
	}

	votesReturned := 0
	votesRequested := len(s.state.peers)
	for {
		select {
		case _, ok := <-s.Quit:
			if !ok {
				Log.Info("Exiting election loop")
				return
			}
		case vote := <-voteChannel:
			if term != s.state.GetCurrentTerm() || s.state.GetCurrentState() != CANDIDATE {
				return
			}
			votesReturned++
			if vote != nil {
				if vote.Term > s.state.GetCurrentTerm() {
					Log.Info("Stopping election, higher term encountered.")
					s.state.SetCurrentTerm(vote.Term)
					s.state.SetCurrentState(FOLLOWER)
					return
				} else {
					if vote.VoteGranted == true {
						votesGranted++
						Log.Info("Vote granted. Current votes :", votesGranted)
						if votesGranted >= votesRequired {
							Log.Info("Node elected leader with", votesGranted, " votes")
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
	nextIndex := s.state.GetNextIndex(node)
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
						s.state.SetNextIndex(node, nextIndex-1)
					}
				} else if response.Success {
					if s.state.GetCurrentState() == LEADER {
						s.state.SetNextIndex(node, nextIndex+1)
						s.state.SetMatchIndex(node, nextIndex)
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
				Log.Info("Exiting leading managment loop")
				return
			}
		case <-s.state.StartLeading:
			Log.Info("Started leading for term ", s.state.GetCurrentTerm())
			s.state.leaderState = newLeaderState(true, &s.state.peers, s.state.log.GetMostRecentIndex())
			for i := 0; i < len(s.state.peers); i++ {
				s.Wait.Add(1)
				go s.sendHeartBeat(&s.state.peers[i])
			}
			timer := time.NewTimer(HEARTBEAT * time.Millisecond)
			for {
				select {
				case _, ok := <-s.Quit:
					if !ok {
						Log.Info("Exiting heartbeat loop")
						return
					}
				case <-s.state.StopLeading:
					break
				case <-s.state.SendAppendEntries:
					timer.Reset(HEARTBEAT * time.Millisecond)
					s.ElectionTimeoutReset <- true
					for i := 0; i < len(s.state.peers); i++ {
						s.Wait.Add(1)
						go s.sendHeartBeat(&s.state.peers[i])
					}
				case <-timer.C:
					timer.Reset(HEARTBEAT * time.Millisecond)
					s.ElectionTimeoutReset <- true
					for i := 0; i < len(s.state.peers); i++ {
						s.Wait.Add(1)
						go s.sendHeartBeat(&s.state.peers[i])
					}
				}
			}
		}
	}
}

func newRaftNetworkServer(nodeId string, peers []Node) *RaftNetworkServer {
	raftServer := &RaftNetworkServer{state: newRaftState(nodeId, peers)}
	raftServer.Wait.Add(3)
	raftServer.ElectionTimeoutReset = make(chan bool, 100)
	raftServer.Quit = make(chan bool)
	go raftServer.electionTimeOut()
	go raftServer.manageElections()
	go raftServer.manageLeading()
	return raftServer
}

func startRaft(lis *net.Listener, nodeId string, peers []Node) (*RaftNetworkServer, *grpc.Server) {
	var opts []grpc.ServerOption
	srv := grpc.NewServer(opts...)
	raftServer := newRaftNetworkServer(nodeId, peers)
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
