package activitylogger

import (
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
	"path"
	"strconv"
)

// listenAppendLog listens to the appendLogChan for incoming log entries.
// ensures synchronus ordered logging to eliminate race conditions
func listenAppendLog() {
	for {
		select {
		case le := <-appendLogChan:
			protoData := getEntryDate(le)
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
		case <-killChan:
			return
		}
	}
}

func getEntryDate(le logEntry) []byte {
	var message proto.Message

	switch le.EntryType {
	case typeChmod:
		path, success := le.Params[0].(string)
		failedConversionCheck(success, "chmod", le.Params)
		mode, success := le.Params[1].(uint32)
		failedConversionCheck(success, "chmod", le.Params)
		message = &pb.ChmodRequest{path, mode}

	case typeCreat:
		path, success := le.Params[0].(string)
		failedConversionCheck(success, "creat", le.Params)
		mode, success := le.Params[1].(uint32)
		failedConversionCheck(success, "creat", le.Params)
		message = &pb.CreatRequest{path, mode}

	case typeLink:
		oldPath, success := le.Params[0].(string)
		failedConversionCheck(success, "link", le.Params)
		newPath, success := le.Params[1].(string)
		failedConversionCheck(success, "link", le.Params)
		message = &pb.LinkRequest{oldPath, newPath}

	case typeMkdir:
		path, success := le.Params[0].(string)
		failedConversionCheck(success, "mkdir", le.Params)
		mode, success := le.Params[1].(uint32)
		failedConversionCheck(success, "mkdir", le.Params)
		message = &pb.MkdirRequest{path, mode}

	case typeRename:
		oldPath, success := le.Params[0].(string)
		failedConversionCheck(success, "rename", le.Params)
		newPath, success := le.Params[1].(string)
		failedConversionCheck(success, "rename", le.Params)
		message = &pb.RenameRequest{oldPath, newPath}

	case typeRmdir:
		path, success := le.Params[0].(string)
		failedConversionCheck(success, "rmdir", le.Params)
		message = &pb.RmdirRequest{path}

	case typeSymLink:
		oldPath, success := le.Params[0].(string)
		failedConversionCheck(success, "symLink", le.Params)
		newPath, success := le.Params[1].(string)
		failedConversionCheck(success, "symLink", le.Params)
		message = &pb.LinkRequest{oldPath, newPath}

	case typeTruncate:
		path, success := le.Params[0].(string)
		failedConversionCheck(success, "truncate", le.Params)
		length, success := le.Params[1].(uint64)
		failedConversionCheck(success, "truncate", le.Params)
		message = &pb.TruncateRequest{path, length}

	case typeUnlink:
		path, success := le.Params[0].(string)
		failedConversionCheck(success, "unlink", le.Params)
		message = &pb.UnlinkRequest{path}

	case typeUtimes:
		//(path string, accessSeconds, accessNanoSeconds, modifySeconds, modifyNanoSeconds int64)
		path, success := le.Params[0].(string)
		failedConversionCheck(success, "utimes", le.Params)
		sec, success := le.Params[1].(int64)
		failedConversionCheck(success, "utimes", le.Params)
		nanSec, success := le.Params[2].(int64)
		failedConversionCheck(success, "utimes", le.Params)
		modSec, success := le.Params[3].(int64)
		failedConversionCheck(success, "utimes", le.Params)
		modNanSec, success := le.Params[4].(int64)
		failedConversionCheck(success, "utimes", le.Params)
		message = &pb.UtimesRequest{path, sec, nanSec, modSec, modNanSec}

	case typeWrite:
		//(path string, data []byte, offset, length uint64)
		path, success := le.Params[0].(string)
		failedConversionCheck(success, "write", le.Params)
		data, success := le.Params[1].([]byte)
		failedConversionCheck(success, "write", le.Params)
		off, success := le.Params[2].(uint64)
		failedConversionCheck(success, "write", le.Params)
		len, success := le.Params[3].(uint64)
		failedConversionCheck(success, "write", le.Params)
		message = &pb.WriteRequest{path, data, off, len}

	default:
		log.Fatalln("Activity logger, unrecognised EntryType:", le.EntryType)
	}

	data, err := proto.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}
	return data
}

// failedConversionCheck is a helper function for error checking for getEntryData
func failedConversionCheck(success bool, logType string, params ...interface{}) {
	if !success {
		log.Fatalln("Activity logger: Bad parameters for", logType, "logEntry\n", params)
	}
}

// Append appends a new entry to the activity log
func AppendLog(entryType uint8, params ...interface{}) {
	appendLogChan <- logEntry{entryType, params}
}
