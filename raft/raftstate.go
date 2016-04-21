package raft

import (
	"encoding/json"
	"fmt"
	"github.com/cpssd/paranoid/pfsd/keyman"
	pb "github.com/cpssd/paranoid/proto/raft"
	"github.com/cpssd/paranoid/raft/raftlog"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

const (
	FOLLOWER int = iota
	CANDIDATE
	LEADER
	INACTIVE
)

const (
	PersistentStateFileName string = "persistentStateFile"
	LogDirectory            string = "raft_logs"
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

type RaftState struct {
	//Used for testing purposes
	specialNumber uint64

	NodeId       string
	pfsDirectory string
	currentState int

	currentTerm uint64
	votedFor    string
	Log         *raftlog.RaftLog
	commitIndex uint64
	lastApplied uint64

	leaderId      string
	Configuration *Configuration

	StartElection     chan bool
	StartLeading      chan bool
	StopLeading       chan bool
	SendAppendEntries chan bool
	ApplyEntries      chan bool
	LeaderElected     chan bool

	snapshotCounter       int
	performingSnapshot    bool
	SendSnapshot          chan Node
	NewSnapshotCreated    chan bool
	SnapshotCounterAtZero chan bool

	waitingForApplied    bool
	EntryApplied         chan *EntryAppliedInfo
	ConfigurationApplied chan *pb.Configuration

	raftInfoDirectory   string
	persistentStateLock sync.Mutex
	stateChangeLock     sync.Mutex
	ApplyEntryLock      sync.Mutex
}

func (s *RaftState) GetCurrentTerm() uint64 {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	return s.currentTerm
}

func (s *RaftState) SetCurrentTerm(x uint64) {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()

	s.votedFor = ""
	s.currentTerm = x
	s.savePersistentState()
}

func (s *RaftState) GetCurrentState() int {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	return s.currentState
}

func (s *RaftState) SetCurrentState(x int) {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()

	if s.currentState == LEADER {
		s.StopLeading <- true
	}
	s.currentState = x
	if x == CANDIDATE {
		s.StartElection <- true
	}
	if x == LEADER {
		s.setLeaderIdUnsafe(s.NodeId)
		s.StartLeading <- true
	}
}

func (s *RaftState) GetPerformingSnapshot() bool {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	return s.performingSnapshot
}

func (s *RaftState) SetPerformingSnapshot(x bool) {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	s.performingSnapshot = x
}

func (s *RaftState) IncrementSnapshotCounter() {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	s.snapshotCounter++
}

func (s *RaftState) DecrementSnapshotCounter() {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	s.snapshotCounter--
	if s.snapshotCounter == 0 {
		s.SnapshotCounterAtZero <- true
	}
}

func (s *RaftState) GetSnapshotCounterValue() int {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	return s.snapshotCounter
}

func (s *RaftState) GetCommitIndex() uint64 {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	return s.commitIndex
}

func (s *RaftState) SetCommitIndex(x uint64) {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	s.commitIndex = x
	s.SendAppendEntries <- true
	s.ApplyEntries <- true
}

//setCommitIndexUnsafe must only be used when the stateChangeLock has already been locked
func (s *RaftState) setCommitIndexUnsafe(x uint64) {
	s.commitIndex = x
	s.SendAppendEntries <- true
	s.ApplyEntries <- true
}

func (s *RaftState) SetWaitingForApplied(x bool) {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	s.waitingForApplied = x
}

func (s *RaftState) GetWaitingForApplied() bool {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	return s.waitingForApplied
}

func (s *RaftState) GetVotedFor() string {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	return s.votedFor
}

func (s *RaftState) SetVotedFor(x string) {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	s.votedFor = x
	s.savePersistentState()
}

func (s *RaftState) GetLeaderId() string {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	return s.leaderId
}

//setLeaderIdUnsafe must only be used when the stateChangeLock has already been locked
func (s *RaftState) setLeaderIdUnsafe(x string) {
	if s.leaderId == "" {
		s.LeaderElected <- true
	}
	s.leaderId = x
}

func (s *RaftState) SetLeaderId(x string) {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()

	if s.leaderId == "" {
		s.LeaderElected <- true
	}
	s.leaderId = x
}

func (s *RaftState) GetLastApplied() uint64 {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	return s.lastApplied
}

func (s *RaftState) SetLastApplied(x uint64) {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	s.lastApplied = x
	s.savePersistentState()
}

//setLastAppliedUnsafe must only be used when the stateChangeLock has already been locked
func (s *RaftState) setLastAppliedUnsafe(x uint64) {
	s.lastApplied = x
	s.savePersistentState()
}

func (s *RaftState) SetSpecialNumber(x uint64) {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	s.specialNumber = x
	s.savePersistentState()
}

func (s *RaftState) GetSpecialNumber() uint64 {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	return s.specialNumber
}

func (s *RaftState) applyLogEntry(logEntry *pb.LogEntry) *StateMachineResult {
	switch logEntry.Entry.Type {
	case pb.Entry_Demo:
		demoCommand := logEntry.Entry.GetDemo()
		if demoCommand == nil {
			Log.Fatal("Error applying Log to state machine")
		}
		s.specialNumber = demoCommand.Number
	case pb.Entry_ConfigurationChange:
		config := logEntry.Entry.GetConfig()
		if config != nil {
			s.ConfigurationApplied <- config
		} else {
			Log.Fatal("Error applying configuration update")
		}
	case pb.Entry_StateMachineCommand:
		libpfsCommand := logEntry.Entry.GetCommand()
		if libpfsCommand == nil {
			Log.Fatal("Error applying Log to state machine")
		}
		if s.pfsDirectory == "" {
			Log.Fatal("PfsDirectory is not set")
		}
		return PerformLibPfsCommand(s.pfsDirectory, libpfsCommand)
	case pb.Entry_KeyStateCommand:
		keyCommand := logEntry.Entry.GetKeyCommand()
		if keyCommand == nil {
			Log.Fatal("Error applying KeyStateCommand to state machine")
		}
		return PerformKSMCommand(keyman.StateMachine, keyCommand)
	}
	return nil
}

//ApplyLogEntries applys all log entries that have been commited but not yet applied
func (s *RaftState) ApplyLogEntries() {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()
	s.ApplyEntryLock.Lock()
	defer s.ApplyEntryLock.Unlock()

	if s.commitIndex > s.lastApplied {
		for i := s.lastApplied + 1; i <= s.commitIndex; i++ {
			LogEntry, err := s.Log.GetLogEntry(i)
			if err != nil {
				Log.Fatal("Unable to get log entry1:", err)
			}
			result := s.applyLogEntry(LogEntry)
			s.setLastAppliedUnsafe(i)
			if s.waitingForApplied {
				s.EntryApplied <- &EntryAppliedInfo{
					Index:  i,
					Result: result,
				}
			}
		}
	}
}

func (s *RaftState) calculateNewCommitIndex() {
	s.stateChangeLock.Lock()
	defer s.stateChangeLock.Unlock()

	newCommitIndex := s.Configuration.CalculateNewCommitIndex(s.commitIndex, s.currentTerm, s.Log)
	if newCommitIndex > s.commitIndex {
		s.setCommitIndexUnsafe(newCommitIndex)
	}
}

type persistentState struct {
	SpecialNumber uint64 `json:"specialnumber"`
	CurrentTerm   uint64 `json:"currentterm"`
	VotedFor      string `json:"votedfor"`
	LastApplied   uint64 `json:"lastapplied"`
}

func (s *RaftState) savePersistentState() {
	s.persistentStateLock.Lock()
	defer s.persistentStateLock.Unlock()

	perState := &persistentState{
		SpecialNumber: s.specialNumber,
		CurrentTerm:   s.currentTerm,
		VotedFor:      s.votedFor,
		LastApplied:   s.lastApplied,
	}

	persistentStateBytes, err := json.Marshal(perState)
	if err != nil {
		Log.Fatal("Error saving persistent state to disk:", err)
	}

	if _, err := os.Stat(s.raftInfoDirectory); os.IsNotExist(err) {
		Log.Fatal("Raft Info Directory does not exist:", err)
	}

	newPeristentFile := path.Join(s.raftInfoDirectory, PersistentStateFileName+"-new")
	err = ioutil.WriteFile(newPeristentFile, persistentStateBytes, 0600)
	if err != nil {
		Log.Fatal("Error writing new persistent state to disk:", err)
	}

	err = os.Rename(newPeristentFile, path.Join(s.raftInfoDirectory, PersistentStateFileName))
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

func newRaftState(myNodeDetails Node, pfsDirectory, raftInfoDirectory string, testConfiguration *StartConfiguration) *RaftState {
	persistentState := getPersistentState(path.Join(raftInfoDirectory, PersistentStateFileName))
	var raftState *RaftState
	if persistentState == nil {
		raftState = &RaftState{
			specialNumber:      0,
			pfsDirectory:       pfsDirectory,
			NodeId:             myNodeDetails.NodeID,
			currentTerm:        0,
			votedFor:           "",
			Log:                raftlog.New(path.Join(raftInfoDirectory, LogDirectory)),
			commitIndex:        0,
			lastApplied:        0,
			leaderId:           "",
			snapshotCounter:    0,
			performingSnapshot: false,
			Configuration:      newConfiguration(raftInfoDirectory, testConfiguration, myNodeDetails, true),
			raftInfoDirectory:  raftInfoDirectory,
		}
	} else {
		raftState = &RaftState{
			specialNumber:      persistentState.SpecialNumber,
			pfsDirectory:       pfsDirectory,
			NodeId:             myNodeDetails.NodeID,
			currentTerm:        persistentState.CurrentTerm,
			votedFor:           persistentState.VotedFor,
			Log:                raftlog.New(path.Join(raftInfoDirectory, LogDirectory)),
			commitIndex:        0,
			lastApplied:        persistentState.LastApplied,
			leaderId:           "",
			snapshotCounter:    0,
			performingSnapshot: false,
			Configuration:      newConfiguration(raftInfoDirectory, testConfiguration, myNodeDetails, true),
			raftInfoDirectory:  raftInfoDirectory,
		}
	}

	raftState.StartElection = make(chan bool, 100)
	raftState.StartLeading = make(chan bool, 100)
	raftState.StopLeading = make(chan bool, 100)
	raftState.SendAppendEntries = make(chan bool, 100)
	raftState.ApplyEntries = make(chan bool, 100)
	raftState.LeaderElected = make(chan bool, 1)
	raftState.EntryApplied = make(chan *EntryAppliedInfo, 100)
	raftState.NewSnapshotCreated = make(chan bool, 100)
	raftState.SnapshotCounterAtZero = make(chan bool, 100)
	raftState.SendSnapshot = make(chan Node, 100)
	raftState.ConfigurationApplied = make(chan *pb.Configuration, 100)
	return raftState
}
