package activitylogger

import (
	"errors"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"path"
	"strconv"
)

func getLogData(index int) []byte {
	data, err := ioutil.ReadFile(path.Join(logDir, strconv.Itoa(index)))
	if err != nil {
		log.Fatalln("Activity logger failed to read log number:", index, "err:", err)
	}
	return data
}

func getLogEntryFromFileData(data []byte) (LogEntry, error) {
	typ := uint8(data[0])
	protoData := data[1:]

	switch typ {
	case TypeChmod:
		ms := &pb.ChmodRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return ChmodProtoToLogEntry(ms), nil
	case TypeCreat:
		ms := &pb.CreatRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return CreatProtoToLogEntry(ms), nil
	case TypeLink:
		ms := &pb.LinkRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return LinkProtoToLogEntry(ms), nil
	case TypeMkdir:
		ms := &pb.MkdirRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return MkdirProtoToLogEntry(ms), nil
	case TypeRename:
		ms := &pb.RenameRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return RenameProtoToLogEntry(ms), nil
	case TypeRmdir:
		ms := &pb.RmdirRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return RmdirProtoToLogEntry(ms), nil
	case TypeSymLink:
		ms := &pb.LinkRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return SymLinkProtoToLogEntry(ms), nil
	case TypeTruncate:
		ms := &pb.TruncateRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return TruncateProtoToLogEntry(ms), nil
	case TypeUnlink:
		ms := &pb.UnlinkRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return UnlinkProtoToLogEntry(ms), nil
	case TypeUtimes:
		ms := &pb.UtimesRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return UtimesProtoToLogEntry(ms), nil
	case TypeWrite:
		ms := &pb.WriteRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return WriteProtoToLogEntry(ms), nil
	default:
		return LogEntry{}, errors.New("Unknown entry type")
	}
}

// Get returns the LogEntry at a given index
func Get(index int) (LogEntry, error) {
	pause()
	defer resume()
	if index < 0 || index >= currentIndex {
		return LogEntry{}, errors.New("Index out of bounds")
	}
	data := getLogData(index)
	le, err := getLogEntryFromFileData(data)
	if err != nil {
		log.Fatalln("The log:", index, "must have been tampered with, err:", err)
	}
	return le, nil
}

// GetAllEntriesSince returns a list of LogEntrys including and after the index given
func GetAllEntriesSince(index int) ([]LogEntry, error) {
	pause()
	defer resume()
	if index < 0 || index >= currentIndex {
		return nil, errors.New("Index out of bounds")
	}
	entries := make([]LogEntry, currentIndex-index)
	entryIndex := 0
	for logIndex := index; logIndex < currentIndex; logIndex++ {
		entry, err := Get(logIndex)
		if err != nil {
			return nil, err
		}
		entries[entryIndex] = entry
		entryIndex++
	}
	return entries, nil
}
