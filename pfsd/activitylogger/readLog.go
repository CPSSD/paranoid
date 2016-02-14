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
	case typeChmod:
		ms := &pb.ChmodRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return ChmodProtoToLogEntry(ms), nil
	case typeCreat:
		ms := &pb.CreatRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return CreatProtoToLogEntry(ms), nil
	case typeLink:
		ms := &pb.LinkRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return LinkProtoToLogEntry(ms), nil
	case typeMkdir:
		ms := &pb.MkdirRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return MkdirProtoToLogEntry(ms), nil
	case typeRename:
		ms := &pb.RenameRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return RenameProtoToLogEntry(ms), nil
	case typeRmdir:
		ms := &pb.RmdirRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return RmdirProtoToLogEntry(ms), nil
	case typeSymLink:
		ms := &pb.LinkRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return SymLinkProtoToLogEntry(ms), nil
	case typeTruncate:
		ms := &pb.TruncateRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return TruncateProtoToLogEntry(ms), nil
	case typeUnlink:
		ms := &pb.UnlinkRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return UnlinkProtoToLogEntry(ms), nil
	case typeUtimes:
		ms := &pb.UtimesRequest{}
		err := proto.Unmarshal(protoData, ms)
		if err != nil {
			log.Fatalln(err)
		}
		return UtimesProtoToLogEntry(ms), nil
	case typeWrite:
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

// Get returns the Logentry at a given index
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
