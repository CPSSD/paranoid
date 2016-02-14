package activitylogger

import (
	"github.com/golang/protobuf/proto"
	"log"
	"os"
	"path"
	"strconv"
)

// listenAppendLog listens to the appendLogChan for incoming log entries.
// ensures synchronus ordered logging to eliminate race conditions
func listenAppendLog() {
	paused = false
	for {
		select {
		case le := <-appendLogChan:
			protoData := getEntryData(le)
			entryTypeData := make([]byte, 1)
			entryTypeData[0] = byte(le.EntryType)
			fileData := append(entryTypeData, protoData...)

			file, err := os.Create(path.Join(logDir, strconv.Itoa(currentIndex)))
			if err != nil {
				log.Fatalln("Activity Logger: failed to create logfile:", currentIndex, "err:", err)
			}

			_, err = file.Write(fileData)
			if err != nil {
				log.Fatalln("Activity Logger: failed to write to logfile:", currentIndex, "err:", err)
			}

			file.Close()
			currentIndex++
		case <-pauseChan:
			paused = true
			<-resumeChan
			paused = false
		case <-killChan:
			return
		}
	}
}

func getEntryData(le logEntry) []byte {
	var message proto.Message

	switch le.EntryType {
	case typeChmod:
		message = LogEntryToChmodProto(le)
	case typeCreat:
		message = LogEntryToCreatProto(le)
	case typeLink:
		message = LogEntryToLinkProto(le)
	case typeMkdir:
		message = LogEntryToMkdirProto(le)
	case typeRename:
		message = LogEntryToRenameProto(le)
	case typeRmdir:
		message = LogEntryToRmdirProto(le)
	case typeSymLink:
		message = LogEntryToSymLinkProto(le)
	case typeTruncate:
		message = LogEntryToTruncateProto(le)
	case typeUnlink:
		message = LogEntryToUnlinkProto(le)
	case typeUtimes:
		message = LogEntryToUtimesProto(le)
	case typeWrite:
		message = LogEntryToWriteProto(le)
	default:
		log.Fatalln("Activity logger, unrecognised EntryType:", le.EntryType)
	}

	data, err := proto.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}
	return data
}
