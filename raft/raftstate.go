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
	IsLeader   bool
	NextIndex  []uint64
	MatchIndex []uint64
}

func newLeaderState() *LeaderState {
	return &LeaderState{
		IsLeader: false,
	}
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
	StopElection  chan bool
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
		//Need to step down as leader. Stop heartbeat loop
	}
	if s.currentState == CANDIDATE {
		//Need to end current election
		s.StopElection <- true
	}
	s.currentState = x
	if x == CANDIDATE {
		//Need to start an eleciton.
		s.StartElection <- true
	}
	if x == LEADER {
		//Need to start election loop
	}
}

func (s *RaftState) GetVotedFor() string {
	return s.votedFor
}

func (s *RaftState) SetVotedFor(x string) {
	s.votedFor = x
}

//Will involve reading from disk in the future
func newRaftState(nodeId string, peers []Node) *RaftState {
	return &RaftState{
		nodeId:      nodeId,
		peers:       peers,
		currentTerm: 0,
		votedFor:    "",
		log:         newRaftLog(),
		commitIndex: 0,
		lastApplied: 0,
		leaderState: newLeaderState(),
	}
}
