package activitylogger

import (
	"errors"
	pb "github.com/cpssd/paranoid/proto/activitylogger"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"path"
	"strconv"
)

// GetEntry will read an entry at the given index returning
// the protobuf and an error if something went wrong
func (al *ActivityLogger) GetEntry(index uint64) (entry *pb.Entry, err error) {
	al.indexLock.Lock()
	defer al.indexLock.Unlock()

	if index < 1 || index >= al.currentIndex {
		return nil, errors.New("Index out of bounds")
	}

	fileIndex := ci2fi(index)
	filePath := path.Join(al.logDir, strconv.FormatUint(fileIndex, 10))
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.New("Failed to read logfile")
	}

	entry = &pb.Entry{}
	err = proto.Unmarshal(fileData, entry)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal file data")
	}

	return entry, nil
}

// GetEntriesSince returns a list of entries including and after the one
// at the given index, and an error object if somethign went wrong
func (al *ActivityLogger) GetEntriesSince(index uint64) (entries []*pb.Entry, err error) {
	al.indexLock.Lock()
	defer al.indexLock.Unlock()

	if index < 1 || index >= al.currentIndex {
		return nil, errors.New("Index out of bounds")
	}

	entries = make([]*pb.Entry, al.currentIndex-index)
	for i := index; i < al.currentIndex; i++ {
		fileIndex := ci2fi(i)
		filePath := path.Join(al.logDir, strconv.FormatUint(fileIndex, 10))
		fileData, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, errors.New("Failed to read logfile")
		}

		var entry *pb.Entry
		err = proto.Unmarshal(fileData, entry)
		if err != nil {
			return nil, errors.New("Failed to Unmarshal file data")
		}

		entries[i-index] = entry
	}

	return entries, nil
}
