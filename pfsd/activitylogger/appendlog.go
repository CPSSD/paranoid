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
		case <-killChan:
			return
		case <-pauseChan:
			paused = true
			<-resumeChan
			paused = false
		case le := <-appendLogChan:
			protoData := getEntryData(le)
			entryTypeData := make([]byte, 1, 1)
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
		}
	}
}

func getEntryData(le LogEntry) []byte {
	var message proto.Message

	switch le.EntryType {
	case TypeChmod:
		message = LogEntryToChmodProto(le)
	case TypeCreat:
		message = LogEntryToCreatProto(le)
	case TypeLink:
		message = LogEntryToLinkProto(le)
	case TypeMkdir:
		message = LogEntryToMkdirProto(le)
	case TypeRename:
		message = LogEntryToRenameProto(le)
	case TypeRmdir:
		message = LogEntryToRmdirProto(le)
	case TypeSymLink:
		message = LogEntryToSymLinkProto(le)
	case TypeTruncate:
		message = LogEntryToTruncateProto(le)
	case TypeUnlink:
		message = LogEntryToUnlinkProto(le)
	case TypeUtimes:
		message = LogEntryToUtimesProto(le)
	case TypeWrite:
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

func pause() {
	pauseChan <- true
}

func resume() {
	resumeChan <- true
}

// Appendlog adds an entry to the log
func Appendlog(typ uint8, params ...interface{}) {
	appendLogChan <- newLogEntry(typ, params...)
}
