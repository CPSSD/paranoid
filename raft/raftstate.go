package raft

import (
	"encoding/json"
	"fmt"
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
	LeaderElected     chan bool

	waitingForApplied    bool
	EntryApplied         chan *EntryAppliedInfo
	ConfigurationApplied chan *pb.Configuration
	ApplyEntriesLock     sync.Mutex

	raftInfoDirectory   string
	persistentStateLock sync.Mutex
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
		s.SetLeaderId(s.NodeId)
		s.StartLeading <- true
	}
}

func (s *RaftState) GetCommitIndex() uint64 {
	return s.commitIndex
}

func (s *RaftState) SetCommitIndex(x uint64) {
	s.commitIndex = x
	s.SendAppendEntries <- true
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

func (s *RaftState) SetSpecialNumber(x uint64) {
	s.specialNumber = x
	s.savePersistentState()
}

func (s *RaftState) GetSpecialNumber() uint64 {
	return s.specialNumber
}

func (s *RaftState) applyLogEntry(logEntry *pb.LogEntry) *StateMachineResult {
	if logEntry.Entry.Type == pb.Entry_Demo {
		demoCommand := logEntry.Entry.GetDemo()
		if demoCommand == nil {
			Log.Fatal("Error applying Log to state machine")
		}
		s.SetSpecialNumber(demoCommand.Number)
	} else if logEntry.Entry.Type == pb.Entry_ConfigurationChange {
		config := logEntry.Entry.GetConfig()
		if config != nil {
			s.ConfigurationApplied <- config
		} else {
			Log.Fatal("Error applying configuration update")
		}
	} else if logEntry.Entry.Type == pb.Entry_StateMachineCommand {
		libpfsCommand := logEntry.Entry.GetCommand()
		if libpfsCommand == nil {
			Log.Fatal("Error applying Log to state machine")
		}
		if s.pfsDirectory == "" {
			Log.Fatal("PfsDirectory is not set")
		}
		return performLibPfsCommand(s.pfsDirectory, libpfsCommand)
	}
	return nil
}

func (s *RaftState) ApplyLogEntries() {
	s.ApplyEntriesLock.Lock()
	defer s.ApplyEntriesLock.Unlock()
	lastApplied := s.GetLastApplied()
	commitIndex := s.GetCommitIndex()
	if commitIndex > lastApplied {
		for i := lastApplied + 1; i <= commitIndex; i++ {
			LogEntry, err := s.Log.GetLogEntry(i)
			if err != nil {
				Log.Fatal("Unable to get log entry1:", err)
			}
			result := s.applyLogEntry(LogEntry)
			s.SetLastApplied(i)
			if s.GetWaitingForApplied() {
				s.EntryApplied <- &EntryAppliedInfo{
					Index:  i,
					Result: result,
				}
			}
		}
	}
}

func (s *RaftState) calculateNewCommitIndex() {
	lastCommitIndex := s.GetCommitIndex()
	currentTerm := s.GetCurrentTerm()
	newCommitIndex := s.Configuration.CalculateNewCommitIndex(lastCommitIndex, currentTerm, s.Log)

	if currentTerm == s.GetCurrentTerm() {
		if newCommitIndex > s.GetCommitIndex() {
			s.SetCommitIndex(newCommitIndex)
			s.ApplyLogEntries()
		}
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
		SpecialNumber: s.GetSpecialNumber(),
		CurrentTerm:   s.GetCurrentTerm(),
		VotedFor:      s.GetVotedFor(),
		LastApplied:   s.GetLastApplied(),
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
			specialNumber:     0,
			pfsDirectory:      pfsDirectory,
			NodeId:            myNodeDetails.NodeID,
			currentTerm:       0,
			votedFor:          "",
			Log:               raftlog.New(path.Join(raftInfoDirectory, LogDirectory)),
			commitIndex:       0,
			lastApplied:       0,
			leaderId:          "",
			Configuration:     newConfiguration(raftInfoDirectory, testConfiguration, myNodeDetails),
			raftInfoDirectory: raftInfoDirectory,
		}
	} else {
		raftState = &RaftState{
			specialNumber:     persistentState.SpecialNumber,
			pfsDirectory:      pfsDirectory,
			NodeId:            myNodeDetails.NodeID,
			currentTerm:       persistentState.CurrentTerm,
			votedFor:          persistentState.VotedFor,
			Log:               raftlog.New(path.Join(raftInfoDirectory, LogDirectory)),
			commitIndex:       0,
			lastApplied:       persistentState.LastApplied,
			leaderId:          "",
			Configuration:     newConfiguration(raftInfoDirectory, testConfiguration, myNodeDetails),
			raftInfoDirectory: raftInfoDirectory,
		}
	}

	raftState.StartElection = make(chan bool, 100)
	raftState.StartLeading = make(chan bool, 100)
	raftState.StopLeading = make(chan bool, 100)
	raftState.SendAppendEntries = make(chan bool, 100)
	raftState.LeaderElected = make(chan bool, 1)
	raftState.EntryApplied = make(chan *EntryAppliedInfo, 100)
	raftState.ConfigurationApplied = make(chan *pb.Configuration, 100)
	return raftState
}
