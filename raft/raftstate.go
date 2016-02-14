package raft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
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

type RaftState struct {
	SpecialNumber uint64

	nodeId       string
	currentState int
	peers        []Node

	currentTerm uint64
	votedFor    string
	log         *RaftLog
	commitIndex uint64
	lastApplied uint64

	leaderId    string
	leaderState *LeaderState

	StartElection     chan bool
	StartLeading      chan bool
	StopLeading       chan bool
	SendAppendEntries chan bool
	LeaderElected     chan bool

	waitingForApplied bool
	EntryApplied      chan uint64
	ApplyEntriesLock  sync.Mutex

	persistentStateFile string
	persistentStateLock sync.Mutex
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
		leaderState.NextIndex[i] = lastLogIndex + 1
		leaderState.MatchIndex[i] = 0
	}
	return leaderState
}

func (s *RaftState) GetNextIndex(node *Node) uint64 {
	for i := 0; i < len(s.peers); i++ {
		if s.peers[i].NodeID == node.NodeID {
			return s.leaderState.NextIndex[i]
		}
	}
	Log.Fatal("Could not get nextIndex. Node not found")
	return 0
}

func (s *RaftState) GetMatchIndex(node *Node) uint64 {
	for i := 0; i < len(s.peers); i++ {
		if s.peers[i].NodeID == node.NodeID {
			return s.leaderState.MatchIndex[i]
		}
	}
	Log.Fatal("Could not get matchIndex. Node not found")
	return 0
}

func (s *RaftState) SetNextIndex(node *Node, x uint64) {
	for i := 0; i < len(s.peers); i++ {
		if s.peers[i].NodeID == node.NodeID {
			s.leaderState.NextIndex[i] = x
			return
		}
	}
	Log.Fatal("Could not set next index")
}

func (s *RaftState) SetMatchIndex(node *Node, x uint64) {
	for i := 0; i < len(s.peers); i++ {
		if s.peers[i].NodeID == node.NodeID {
			s.leaderState.MatchIndex[i] = x
			return
		}
	}
	Log.Fatal("Could not set match index")
}

func (s *RaftState) GetCurrentTerm() uint64 {
	return s.currentTerm
}

func (s *RaftState) SetCurrentTerm(x uint64) {
	s.votedFor = ""
	s.currentTerm = x
	s.savePersistentState()
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
		s.SetLeaderId(s.nodeId)
		s.StartLeading <- true
	}
}

func (s *RaftState) GetCommitIndex() uint64 {
	return s.commitIndex
}

func (s *RaftState) SetCommitIndex(x uint64) {
	s.commitIndex = x
}

func (s *RaftState) SetWaitingForApplied(x bool) {
	s.waitingForApplied = x
}

func (s *RaftState) GetWaitingForApplied() bool {
	return s.waitingForApplied
}

func (s *RaftState) GetVotedFor() string {
	return s.votedFor
}

func (s *RaftState) SetVotedFor(x string) {
	s.votedFor = x
	s.savePersistentState()
}

func (s *RaftState) GetLeaderId() string {
	return s.leaderId
}

func (s *RaftState) SetLeaderId(x string) {
	if s.leaderId == "" {
		s.LeaderElected <- true
	}
	s.leaderId = x
}

func (s *RaftState) GetLastApplied() uint64 {
	return s.lastApplied
}

func (s *RaftState) SetLastApplied(x uint64) {
	s.lastApplied = x
	s.savePersistentState()
}

func (s *RaftState) ApplyLogEntries() {
	s.ApplyEntriesLock.Lock()
	defer s.ApplyEntriesLock.Unlock()
	lastApplied := s.GetLastApplied()
	commitIndex := s.GetCommitIndex()
	if commitIndex > lastApplied {
		for i := lastApplied + 1; i <= commitIndex; i++ {
			logEntry := s.log.GetLogEntry(i)
			s.SpecialNumber = logEntry.Entry.Number
			s.SetLastApplied(i)
			if s.GetWaitingForApplied() {
				s.EntryApplied <- i
			}
		}
	}
}

func (s *RaftState) calculateNewCommitIndex() {
	majority := len(s.peers) / 2
	if len(s.peers)%2 == 0 {
		majority++
	}
	lastCommitIndex := s.GetCommitIndex()
	for i := lastCommitIndex + 1; i <= s.log.GetMostRecentIndex(); i++ {
		if s.log.GetLogEntry(i).Term == s.GetCurrentTerm() {
			count := 1
			for j := 0; j < len(s.leaderState.MatchIndex); j++ {
				if s.leaderState.MatchIndex[j] >= i {
					count++
				}
			}
			if count < majority {
				break
			}
			s.SetCommitIndex(i)
		}
	}
	if s.GetCommitIndex() > lastCommitIndex {
		s.ApplyLogEntries()
	}
}

type persistentState struct {
	CurrentTerm uint64 `json:"currentterm"`
	VotedFor    string `json:"votedfor"`
	LastApplied uint64 `json:"lastapplied"`
}

func (s *RaftState) savePersistentState() {
	s.persistentStateLock.Lock()
	defer s.persistentStateLock.Unlock()
	perState := &persistentState{
		CurrentTerm: s.GetCurrentTerm(),
		VotedFor:    s.GetVotedFor(),
		LastApplied: s.GetLastApplied(),
	}

	persistentStateBytes, err := json.Marshal(perState)
	if err != nil {
		Log.Fatal("Error saving persistent state to disk:", err)
	}

	err = ioutil.WriteFile(s.persistentStateFile, persistentStateBytes, 0600)
	if err != nil {
		Log.Fatal("Error saving persistent state to disk:", err)
	}
}

func getPersistentState(persistentStateFile string) *persistentState {
	if _, err := os.Stat(persistentStateFile); os.IsNotExist(err) {
		return nil
	}
	persistentFileContents, err := ioutil.ReadFile(persistentStateFile)
	if err != nil {
		Log.Fatal("Error reading persistent state from disk:", err)
	}

	perState := &persistentState{}
	err = json.Unmarshal(persistentFileContents, &perState)
	if err != nil {
		Log.Fatal("Error reading persistent state from disk:", err)
	}
	return perState
}

func newRaftState(nodeId, persistentStateFile string, peers []Node) *RaftState {
	persistentState := getPersistentState(persistentStateFile)
	var raftState *RaftState
	if persistentState == nil {
		raftState = &RaftState{
			nodeId:              nodeId,
			peers:               peers,
			currentTerm:         0,
			votedFor:            "",
			log:                 newRaftLog(),
			commitIndex:         0,
			lastApplied:         0,
			leaderId:            "",
			leaderState:         newLeaderState(false, nil, 0),
			persistentStateFile: persistentStateFile,
		}
	} else {
		raftState = &RaftState{
			nodeId:              nodeId,
			peers:               peers,
			currentTerm:         persistentState.CurrentTerm,
			votedFor:            persistentState.VotedFor,
			log:                 newRaftLog(),
			commitIndex:         0,
			lastApplied:         persistentState.LastApplied,
			leaderId:            "",
			leaderState:         newLeaderState(false, nil, 0),
			persistentStateFile: persistentStateFile,
		}
	}

	raftState.StartElection = make(chan bool, 100)
	raftState.StartLeading = make(chan bool, 100)
	raftState.StopLeading = make(chan bool, 100)
	raftState.SendAppendEntries = make(chan bool, 100)
	raftState.LeaderElected = make(chan bool, 1)
	raftState.EntryApplied = make(chan uint64, 100)
	return raftState
}
