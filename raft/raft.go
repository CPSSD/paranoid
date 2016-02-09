package raft

import (
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
	ELECTION_TIMEOUT int = 300
)

var (
	Log *logger.ParanoidLogger
)

type RaftNetworkServer struct {
	state                *RaftState
	Wait                 sync.WaitGroup
	Quit                 chan bool
	ElectionTimeoutReset chan bool
}

func (s *RaftNetworkServer) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	if req.Term < s.state.GetCurrentTerm() {
		return &pb.AppendEntriesResponse{s.state.GetCurrentTerm(), false}, nil
	}

	if req.Term > s.state.GetCurrentTerm() {
		s.state.SetCurrentTerm(req.Term)
		s.state.SetCurrentState(FOLLOWER)
	}

	preLogEntry := s.state.log.GetLogEntry(req.PrevLogIndex)
	if preLogEntry == nil || preLogEntry.Term != req.PrevLogTerm {
		return &pb.AppendEntriesResponse{s.state.GetCurrentTerm(), false}, nil
	}

	return &pb.AppendEntriesResponse{s.state.GetCurrentTerm(), false}, nil
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

func getRandomElectionTimeout() time.Duration {
	return time.Duration(ELECTION_TIMEOUT + rand.Intn(ELECTION_TIMEOUT))
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

func Dial(node Node) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTimeout(2*time.Second))
	//TODO: tls support
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(node.String(), opts...)
	return conn, err
}

func (s *RaftNetworkServer) requestPeerVote(node Node, term uint64, voteChannel chan *pb.RequestVoteResponse) {
	defer s.Wait.Done()
	for {
		if term != s.state.GetCurrentTerm() || s.state.GetCurrentState() != CANDIDATE {
			voteChannel <- nil
			return
		}
		conn, err := Dial(node)
		if err != nil {
			client := pb.NewRaftNetworkClient(conn)
			response, err := client.RequestVote(context.Background(), &pb.RequestVoteRequest{s.state.GetCurrentTerm(),
				s.state.nodeId,
				s.state.log.GetMostRecentIndex(),
				s.state.log.GetMostRecentTerm()})
			if err != nil {
				voteChannel <- response
				return
			}
		}
	}
}

func (s *RaftNetworkServer) getRequiredVotes() int {
	return 1
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

	voteChannel := make(chan *pb.RequestVoteResponse)
	for i := 0; i < len(s.state.peers); i++ {
		s.Wait.Add(1)
		go s.requestPeerVote(s.state.peers[i], term, voteChannel)
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
		case <-s.state.StopElection:
			Log.Info("Election ended")
			return
		case vote := <-voteChannel:
			if term != s.state.GetCurrentTerm() || s.state.GetCurrentState() != CANDIDATE {
				return
			}
			votesReturned++
			if vote != nil {
				if vote.Term > s.state.GetCurrentTerm() {
					s.state.SetCurrentTerm(vote.Term)
					s.state.SetCurrentState(FOLLOWER)
				} else {
					if vote.VoteGranted == true {
						votesGranted++
						if votesGranted >= votesRequired {
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

func newRaftNetworkServer(nodeId string, peers []Node) *RaftNetworkServer {
	return &RaftNetworkServer{state: newRaftState(nodeId, peers)}
}

func startRaft(lis *net.Listener) {
	var opts []grpc.ServerOption
	srv := grpc.NewServer(opts...)
	raftServer := newRaftNetworkServer("Hi", []Node{})
	pb.RegisterRaftNetworkServer(srv, raftServer)
	srv.Serve(*lis)
}
