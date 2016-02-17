package activitylogger

import (
	"errors"
	pb "github.com/cpssd/paranoid/proto/activitylogger"
	"github.com/golang/protobuf/proto"
	"os"
	"path"
	"strconv"
)

// WriteEntry will write the entry provided and return the
// index of the entry and an error object if somethign went wrong
func WriteEntry(en *pb.Entry) (index int, err error) {
	indexLock.Lock()
	defer indexLock.Unlock()

	fileIndex := currentIndex
	filePath := path.Join(logDir, strconv.Itoa(fileIndex))

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
			pLog.Fatal("Failed to write proto to file at index: ", fileIndex,
				" and received an erro when trying to remove the created logfile, err: ", err1)
		}
		return 0, errors.New("Failed to write proto to file")
	}

	currentIndex++
	return fileIndex, nil
}
