package raft

import (
	"fmt"
)

const (
	FOLLOWER int = iota
	CANDIDATE
	LEADER
)

type Node struct {
	IP         string
	Port       string
	CommonName string
	NodeID     string
}

func (n Node) String() string {
	return fmt.Sprintf("%s:%s", n.IP, n.Port)
}

type LeaderState struct {
	NextIndex  []uint64
	MatchIndex []uint64
}

func newLeaderState(isLeader bool, peers *[]Node, lastLogIndex uint64) *LeaderState {
	if isLeader == false {
		return &LeaderState{
			NextIndex:  make([]uint64, 0),
			MatchIndex: make([]uint64, 0),
		}
	}
	leaderState := &LeaderState{
		NextIndex:  make([]uint64, len(*peers)),
		MatchIndex: make([]uint64, len(*peers)),
	}
	for i := 0; i < len(*peers); i++ {
		leaderState.NextIndex[i] = lastLogIndex
		leaderState.MatchIndex[i] = 0
	}
	return leaderState
}

func (s *RaftState) GetNextIndex(node Node) uint64 {
	for i := 0; i < len(s.peers); i++ {
		if s.peers[i].NodeID == node.NodeID {
			return s.leaderState.NextIndex[i]
		}
	}
	Log.Fatal("Could not get nextIndex. Node not found")
	return 0
}

func (s *RaftState) GetMatchIndex(node Node) uint64 {
	for i := 0; i < len(s.peers); i++ {
		if s.peers[i].NodeID == node.NodeID {
			return s.leaderState.MatchIndex[i]
		}
	}
	Log.Fatal("Could not get matchIndex. Node not found")
	return 0
}

func (s *RaftState) SetNextIndex(node Node, x uint64) {
	for i := 0; i < len(s.peers); i++ {
		if s.peers[i].NodeID == node.NodeID {
			s.leaderState.NextIndex[i] = x
			return
		}
	}
	Log.Fatal("Could not set next index")
}

func (s *RaftState) SetMatchIndex(node Node, x uint64) {
	for i := 0; i < len(s.peers); i++ {
		if s.peers[i].NodeID == node.NodeID {
			s.leaderState.MatchIndex[i] = x
			return
		}
	}
	Log.Fatal("Could not set match index")
}

type RaftState struct {
	nodeId       string
	currentState int
	peers        []Node

	currentTerm uint64
	votedFor    string
	log         *RaftLog
	commitIndex uint64
	lastApplied uint64

	leaderState *LeaderState

	StartElection chan bool
	StartLeading  chan bool
	StopLeading   chan bool
}

func (s *RaftState) GetCurrentTerm() uint64 {
	return s.currentTerm
}

func (s *RaftState) SetCurrentTerm(x uint64) {
	s.votedFor = ""
	s.currentTerm = x
}

func (s *RaftState) GetCurrentState() int {
	return s.currentState
}

func (s *RaftState) SetCurrentState(x int) {
	if s.currentState == LEADER {
		s.StopLeading <- true
	}
	s.currentState = x
	if x == CANDIDATE {
		s.StartElection <- true
	}
	if x == LEADER {
		s.StartLeading <- true
	}
}

func (s *RaftState) GetCommitIndex() uint64 {
	return s.commitIndex
}

func (s *RaftState) SetCommitIndex(x uint64) {
	s.commitIndex = x
}

func (s *RaftState) GetVotedFor() string {
	return s.votedFor
}

func (s *RaftState) SetVotedFor(x string) {
	s.votedFor = x
}

//Will involve reading from disk in the future
func newRaftState(nodeId string, peers []Node) *RaftState {
	raftState := &RaftState{
		nodeId:      nodeId,
		peers:       peers,
		currentTerm: 0,
		votedFor:    "",
		log:         newRaftLog(),
		commitIndex: 0,
		lastApplied: 0,
		leaderState: newLeaderState(false, nil, 0),
	}
	raftState.StartElection = make(chan bool, 100)
	raftState.StartLeading = make(chan bool, 100)
	raftState.StopLeading = make(chan bool, 100)
	return raftState
}
