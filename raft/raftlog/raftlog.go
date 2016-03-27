package raftlog

import (
	"encoding/json"
	"github.com/cpssd/paranoid/logger"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

const (
	LogEntryDirectoryName string = "log_entries"
	LogMetaFileName       string = "logmetainfo"
)

var Log *logger.ParanoidLogger

type PeristentLogState struct {
	LogSizeBytes uint64 `json:"logsizebytes"`
	StartIndex   uint64 `json:"startindex"`
	StartTerm    uint64 `json:"startterm"`
}

// RaftLog is the structure through which logging functinality can be accessed
type RaftLog struct {
	logDir         string
	startIndex     uint64
	startTerm      uint64
	logSizeBytes   uint64
	currentIndex   uint64
	mostRecentTerm uint64
	indexLock      sync.Mutex
}

// New returns an initiated instance of RaftLog
func New(logDirectory string) *RaftLog {
	rl := &RaftLog{
		logDir: logDirectory,
	}

	logEntryDirectory := path.Join(rl.logDir, LogEntryDirectoryName)
	logMetaFile := path.Join(rl.logDir, LogMetaFileName)

	fileDescriptors, err := ioutil.ReadDir(logEntryDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(logEntryDirectory, 0700)
			if err != nil {
				Log.Fatal("failed to create log directory:", err)
			}
		} else if os.IsPermission(err) {
			Log.Fatal("raft logger does not have permissions for:", logEntryDirectory)
		} else {
			Log.Fatal("unable to read log directory:", err)
		}
	}

	rl.startIndex = 0
	rl.startTerm = 0
	rl.logSizeBytes = 0
	rl.mostRecentTerm = 0
	metaFileContents, err := ioutil.ReadFile(logMetaFile)
	if err != nil {
		if !os.IsNotExist(err) {
			Log.Fatal("unable to read raft log meta information:", err)
		}
	} else {
		metaInfo := &PeristentLogState{}
		err = json.Unmarshal(metaFileContents, metaInfo)
		if err != nil {
			Log.Fatal("unable to read raft log meta information:", err)
		}
		rl.startIndex = metaInfo.StartIndex
		rl.logSizeBytes = metaInfo.LogSizeBytes
		rl.startTerm = metaInfo.StartTerm
		rl.mostRecentTerm = metaInfo.StartTerm
	}

	rl.currentIndex = uint64(len(fileDescriptors)) + rl.startIndex + 1
	if rl.currentIndex > rl.startIndex+1 {
		logEntry, err := rl.GetLogEntry(rl.currentIndex - 1)
		if err != nil {
			Log.Fatal("failed setting up raft logger, could not get most recent term:", err)
		}
		rl.mostRecentTerm = logEntry.Term
	}
	return rl
}

func (rl *RaftLog) saveMetaInfo() {
	metaInfo := &PeristentLogState{
		LogSizeBytes: rl.logSizeBytes,
		StartIndex:   rl.startIndex,
		StartTerm:    rl.startTerm,
	}

	metaInfoJson, err := json.Marshal(metaInfo)
	if err != nil {
		Log.Fatal("unable to save raft log meta information:", err)
	}

	err = ioutil.WriteFile(path.Join(rl.logDir, LogMetaFileName), metaInfoJson, 0600)
	if err != nil {
		Log.Fatal("unable to save raft log meta information:", err)
	}
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

func (rl *RaftLog) setLogSizeBytes(x uint64) {
	rl.logSizeBytes = x
	rl.saveMetaInfo()
}

func (rl *RaftLog) setStartIndex(x uint64) {
	rl.startIndex = x
	rl.saveMetaInfo()
}

func (rl *RaftLog) setStartTerm(x uint64) {
	rl.startTerm = x
	rl.saveMetaInfo()
}
