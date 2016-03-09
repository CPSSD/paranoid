package raftlog

import (
	"github.com/cpssd/paranoid/logger"
	"io/ioutil"
	"os"
	"sync"
)

// RaftLog is the structure through which logging functinality can be accessed
type RaftLog struct {
	logDir         string
	currentIndex   uint64
	mostRecentTerm uint64
	indexLock      sync.Mutex
	pLog           *logger.ParanoidLogger
}

// New returns an initiated instance of RaftLog
func New(logDirectory string) *RaftLog {
	rl := &RaftLog{
		logDir: logDirectory,
		pLog:   logger.New("Raft Logger", "pfsd", logDirectory),
	}
	fileDescriptors, err := ioutil.ReadDir(rl.logDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(rl.logDir, 0700)
			if err != nil {
				rl.pLog.Fatal("failed to create log directory:", err)
			}
		} else if os.IsPermission(err) {
			rl.pLog.Fatal("Raft logger does not have permissions for:", rl.logDir)
		} else {
			rl.pLog.Fatal("Unable to read log directory:", err)
		}
	}
	rl.currentIndex = uint64(len(fileDescriptors) + 1)
	if rl.currentIndex > 1 {
		logEntry, err := rl.GetLogEntry(rl.currentIndex - 1)
		if err != nil {
			rl.pLog.Fatal("Failed to set up raft logger, could not get most recent term:", err)
		}
		rl.mostRecentTerm = logEntry.Term
	} else {
		rl.mostRecentTerm = 0
	}
	return rl
}

// GetMostRecentIndex returns the index of the last log entry
func (rl *RaftLog) GetMostRecentIndex() uint64 {
	rl.indexLock.Lock()
	defer rl.indexLock.Unlock()
	return rl.currentIndex - 1
}

// GetMostRecentTerm returns the term of the last log entry
func (rl *RaftLog) GetMostRecentTerm() uint64 {
	rl.indexLock.Lock()
	defer rl.indexLock.Unlock()
	return rl.mostRecentTerm
}

// fileIndexToStorageIndex converts a fileIndex to a the index it is stored at
func fileIndexToStorageIndex(i uint64) uint64 {
	return i - 1000000
}

// storageIndexToFileIndex converts a storage index to a fileIndex
func storageIndexToFileIndex(i uint64) uint64 {
	return i + 1000000
}
