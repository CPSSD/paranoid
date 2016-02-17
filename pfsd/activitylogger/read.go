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
// the protobuf and an error if somethign went wrong
func GetEntry(index int) (entry *pb.Entry, err error) {
	indexLock.Lock()
	defer indexLock.Unlock()

	if index < 0 || index >= currentIndex {
		return nil, errors.New("Index out of bounds")
	}

	filePath := path.Join(logDir, strconv.Itoa(index))
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
func GetEntriesSince(index int) (entries []*pb.Entry, err error) {
	indexLock.Lock()
	defer indexLock.Unlock()

	if index < 0 || index >= currentIndex {
		return nil, errors.New("Index out of bounds")
	}

	entries = make([]*pb.Entry, currentIndex-index)
	for i := index; i < currentIndex; i++ {
		filePath := path.Join(logDir, strconv.Itoa(i))
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

	return entries, err
}
