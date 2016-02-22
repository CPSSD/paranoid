package activitylogger

import (
	"errors"
	pb "github.com/cpssd/paranoid/proto/raft"
	"github.com/golang/protobuf/proto"
	"os"
	"path"
	"strconv"
)

// WriteEntry will write the entry provided and return the
// index of the entry and an error object if somethign went wrong
func (al *ActivityLogger) WriteEntry(en *pb.LogEntry) (index uint64, err error) {
	al.indexLock.Lock()
	defer al.indexLock.Unlock()

	fileIndex := ci2fi(al.currentIndex)
	filePath := path.Join(al.logDir, strconv.FormatUint(fileIndex, 10))

	protoData, err := proto.Marshal(en)
	if err != nil {
		return 0, errors.New("Failed to Marshal entry")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return 0, errors.New("Unable to create logfile")
	}
	defer file.Close()

	_, err = file.Write(protoData)
	if err != nil {
		err1 := os.Remove(filePath)
		if err1 != nil {
			al.pLog.Fatal("Failed to write proto to file at index: ", fileIndex,
				" and received an erro when trying to remove the created logfile, err: ", err1)
		}
		return 0, errors.New("Failed to write proto to file")
	}

	al.currentIndex++
	return fi2ci(fileIndex), nil
}
