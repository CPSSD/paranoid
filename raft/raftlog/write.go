package raftlog

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/encryption"
	pb "github.com/cpssd/paranoid/proto/raft"
	"github.com/golang/protobuf/proto"
	"os"
	"path"
	"strconv"
)

// AppendEntry will write the entry provided and return the
// index of the entry and an error object if somethign went wrong
func (rl *RaftLog) AppendEntry(en *pb.LogEntry) (index uint64, err error) {
	rl.indexLock.Lock()
	defer rl.indexLock.Unlock()

	fileIndex := storageIndexToFileIndex(rl.currentIndex)
	filePath := path.Join(rl.logDir, strconv.FormatUint(fileIndex, 10))

	protoData, err := proto.Marshal(en)
	if err != nil {
		return 0, errors.New("Failed to Marshal entry")
	}

	if encryption.Encrypted {
		blockSize := encryption.GetCipherSize()
		extraBytes := make([]byte, blockSize-len(protoData)%blockSize)
		protoData = append(protoData, extraBytes...)
		err = encryption.Encrypt(protoData)
		if err != nil {
			return 0, fmt.Errorf("error encrypting log entry: %s", err)
		}
		protoData = append(protoData, byte(len(extraBytes)))
	}

	file, err := os.Create(filePath)
	if err != nil {
		return 0, errors.New("Unable to create logfile")
	}
	defer file.Close()

	_, err = file.Write(protoData)
	if err != nil {
		Log.Error("unable to write proto to file at index", fileIndex)
		err := os.Remove(filePath)
		if err != nil {
			Log.Fatal("unable to remove the created logfile:", err)
		}
		return 0, errors.New("Failed to write proto to file")
	}

	rl.mostRecentTerm = en.Term
	rl.currentIndex++
	return fileIndexToStorageIndex(fileIndex), nil
}
