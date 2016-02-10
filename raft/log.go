package raft

import (
	pb "github.com/cpssd/paranoid/proto/raft"
)

type RaftLog struct {
	logEntries []LogEntry
	startIndex uint64
}

type LogEntry struct {
	Entry pb.Entry
	Term  uint64
}

//Will involve disk reading in the future and possibly some form of caching
func (l *RaftLog) GetLogEntry(index uint64) *LogEntry {
	adjustedIndex := int(index) - int(l.startIndex) - 1
	if adjustedIndex > len(l.logEntries) || adjustedIndex < 0 {
		return nil
	}
	return &l.logEntries[adjustedIndex]
}

func (l *RaftLog) GetMostRecentTerm() uint64 {
	if len(l.logEntries) == 0 {
		return 0
	}
	return l.logEntries[len(l.logEntries)-1].Term
}

func (l *RaftLog) GetMostRecentIndex() uint64 {
	return l.startIndex + uint64(len(l.logEntries))
}

//Will involve reading from disk in the future
func newRaftLog() *RaftLog {
	var logEntries []LogEntry
	var startIndex uint64
	startIndex = 0
	return &RaftLog{logEntries, startIndex}
}
