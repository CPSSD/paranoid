package raftlog

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/encryption"
	pb "github.com/cpssd/paranoid/proto/raft"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"math"
	"path"
	"strconv"
)

func (rl *RaftLog) GetLogEntryUnsafe(index uint64) (entry *pb.LogEntry, err error) {
	if index < 1 || index >= rl.currentIndex {
		return nil, errors.New("index out of bounds")
	}

	if index <= rl.startIndex {
		return nil, ErrIndexBelowStartIndex
	}

	fileIndex := storageIndexToFileIndex(index)
	filePath := path.Join(rl.logDir, LogEntryDirectoryName, strconv.FormatUint(fileIndex, 10))
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("failed to read logfile")
	}

	if encryption.Encrypted {
		err = encryption.Decrypt(fileData[:len(fileData)-1])
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt logfile: %s", err)
		}
		extraEndBytes := int(fileData[len(fileData)-1]) + 1
		fileData = fileData[:len(fileData)-extraEndBytes]
	}

	entry = &pb.LogEntry{}
	err = proto.Unmarshal(fileData, entry)
	if err != nil {
		return nil, errors.New("failed to Unmarshal file data")
	}

	return entry, nil
}

// GetLogEntry will read an entry at the given index returning
// the protobuf and an error if something went wrong
func (rl *RaftLog) GetLogEntry(index uint64) (entry *pb.LogEntry, err error) {
	rl.indexLock.Lock()
	defer rl.indexLock.Unlock()

	return rl.GetLogEntryUnsafe(index)
}

// GetEntriesSince returns a list of entries including and after the one
// at the given index, and an error object if something went wrong
func (rl *RaftLog) GetEntriesSince(index uint64) (entries []*pb.Entry, err error) {
	return rl.GetLogEntries(index, math.MaxUint64-index)
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

// GetLogEntries gets maxCount log entries starting from index. If there are less
// entries than maxCount it gets all of them until the end.
func (rl *RaftLog) GetLogEntries(index, maxCount uint64) (entries []*pb.Entry, err error) {
	rl.indexLock.Lock()
	defer rl.indexLock.Unlock()

	if index < 1 || index >= rl.currentIndex {
		return nil, errors.New("index out of bounds")
	}

	if index <= rl.startIndex {
		return nil, ErrIndexBelowStartIndex
	}

	entries = make([]*pb.Entry, min(rl.currentIndex-index, maxCount))
	for i := index; i < min(rl.currentIndex, index+maxCount); i++ {
		fileIndex := storageIndexToFileIndex(i)
		filePath := path.Join(rl.logDir, LogEntryDirectoryName, strconv.FormatUint(fileIndex, 10))
		fileData, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, errors.New("failed to read logfile")
		}

		if encryption.Encrypted {
			err = encryption.Decrypt(fileData[:len(fileData)-1])
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt logfile: %s", err)
			}
			extraEndBytes := int(fileData[len(fileData)-1]) + 1
			fileData = fileData[:len(fileData)-extraEndBytes]
		}

		entry := &pb.LogEntry{}
		err = proto.Unmarshal(fileData, entry)
		if err != nil {
			return nil, errors.New("failed to unmarshal file data")
		}

		entries[i-index] = entry.Entry
	}

	return entries, nil
}
