package activitylogger

import (
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"log"
)

// LogEntryToChmodProto Converys a LogEntry to a *pb.ChmodRequest
func LogEntryToChmodProto(le LogEntry) *pb.ChmodRequest {
	path, success := le.Params[0].(string)
	failedConversionCheck(success, "chmod", le.Params)
	mode, success := le.Params[1].(uint32)
	failedConversionCheck(success, "chmod", le.Params)
	return &pb.ChmodRequest{path, mode}
}

// LogEntryToCreatProto Converys a LogEntry to a *pb.CreatRequest
func LogEntryToCreatProto(le LogEntry) *pb.CreatRequest {
	path, success := le.Params[0].(string)
	failedConversionCheck(success, "creat", le.Params)
	mode, success := le.Params[1].(uint32)
	failedConversionCheck(success, "creat", le.Params)
	return &pb.CreatRequest{path, mode}
}

// LogEntryToLinkProto Converys a LogEntry to a *pb.LinkRequest
func LogEntryToLinkProto(le LogEntry) *pb.LinkRequest {
	oldPath, success := le.Params[0].(string)
	failedConversionCheck(success, "link", le.Params)
	newPath, success := le.Params[1].(string)
	failedConversionCheck(success, "link", le.Params)
	return &pb.LinkRequest{oldPath, newPath}
}

// LogEntryToMkdirProto Converys a LogEntry to a *pb.MkdirRequest
func LogEntryToMkdirProto(le LogEntry) *pb.MkdirRequest {
	path, success := le.Params[0].(string)
	failedConversionCheck(success, "mkdir", le.Params)
	mode, success := le.Params[1].(uint32)
	failedConversionCheck(success, "mkdir", le.Params)
	return &pb.MkdirRequest{path, mode}
}

// LogEntryToRenameProto Converys a LogEntry to a *pb.RenameRequest
func LogEntryToRenameProto(le LogEntry) *pb.RenameRequest {
	oldPath, success := le.Params[0].(string)
	failedConversionCheck(success, "rename", le.Params)
	newPath, success := le.Params[1].(string)
	failedConversionCheck(success, "rename", le.Params)
	return &pb.RenameRequest{oldPath, newPath}
}

// LogEntryToRmdirProto Converys a LogEntry to a *pb.RmdirRequest
func LogEntryToRmdirProto(le LogEntry) *pb.RmdirRequest {
	path, success := le.Params[0].(string)
	failedConversionCheck(success, "rmdir", le.Params)
	return &pb.RmdirRequest{path}
}

// LogEntryToSymLinkProto Converys a LogEntry to a *pb.LinkRequest
func LogEntryToSymLinkProto(le LogEntry) *pb.LinkRequest {
	oldPath, success := le.Params[0].(string)
	failedConversionCheck(success, "symLink", le.Params)
	newPath, success := le.Params[1].(string)
	failedConversionCheck(success, "symLink", le.Params)
	return &pb.LinkRequest{oldPath, newPath}
}

// LogEntryToTruncateProto Converys a LogEntry to a *pb.TruncateRequest
func LogEntryToTruncateProto(le LogEntry) *pb.TruncateRequest {
	path, success := le.Params[0].(string)
	failedConversionCheck(success, "truncate", le.Params)
	length, success := le.Params[1].(uint64)
	failedConversionCheck(success, "truncate", le.Params)
	return &pb.TruncateRequest{path, length}
}

// LogEntryToUnlinkProto Converys a LogEntry to a *pb.UnlinkRequest
func LogEntryToUnlinkProto(le LogEntry) *pb.UnlinkRequest {
	path, success := le.Params[0].(string)
	failedConversionCheck(success, "unlink", le.Params)
	return &pb.UnlinkRequest{path}
}

// LogEntryToUtimesProto Converys a LogEntry to a *pb.UtimesRequest
func LogEntryToUtimesProto(le LogEntry) *pb.UtimesRequest {
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
	return &pb.UtimesRequest{path, sec, nanSec, modSec, modNanSec}
}

// LogEntryToWriteProto Converys a LogEntry to a *pb.WriteRequest
func LogEntryToWriteProto(le LogEntry) *pb.WriteRequest {
	path, success := le.Params[0].(string)
	failedConversionCheck(success, "write", le.Params)
	data, success := le.Params[1].([]byte)
	failedConversionCheck(success, "write", le.Params)
	off, success := le.Params[2].(uint64)
	failedConversionCheck(success, "write", le.Params)
	len, success := le.Params[3].(uint64)
	failedConversionCheck(success, "write", le.Params)
	return &pb.WriteRequest{path, data, off, len}
}

// ChmodProtoToLogEntry converts a ChmodRequest protobuf into a LogEntry
func ChmodProtoToLogEntry(pro *pb.ChmodRequest) LogEntry {
	return newLogEntry(typeChmod, pro.Path, pro.Mode)
}

// CreatProtoToLogEntry converts a CreatRequest protobuf into a LogEntry
func CreatProtoToLogEntry(pro *pb.CreatRequest) LogEntry {
	return newLogEntry(typeCreat, pro.Path, pro.Permissions)
}

// LinkProtoToLogEntry converts a LinkRequest protobuf into a LogEntry
func LinkProtoToLogEntry(pro *pb.LinkRequest) LogEntry {
	return newLogEntry(typeLink, pro.OldPath, pro.NewPath)
}

// MkdirProtoToLogEntry converts a MkdirRequest protobuf into a LogEntry
func MkdirProtoToLogEntry(pro *pb.MkdirRequest) LogEntry {
	return newLogEntry(typeMkdir, pro.Directory, pro.Mode)
}

// RenameProtoToLogEntry converts a RenameRequest protobuf into a LogEntry
func RenameProtoToLogEntry(pro *pb.RenameRequest) LogEntry {
	return newLogEntry(typeRename, pro.OldPath, pro.NewPath)
}

// RmdirProtoToLogEntry converts a RmdirRequest protobuf into a LogEntry
func RmdirProtoToLogEntry(pro *pb.RmdirRequest) LogEntry {
	return newLogEntry(typeRmdir, pro.Directory)
}

// SymLinkProtoToLogEntry converts a LinkRequest protobuf into a LogEntry
func SymLinkProtoToLogEntry(pro *pb.LinkRequest) LogEntry {
	return newLogEntry(typeSymLink, pro.OldPath, pro.NewPath)
}

// TruncateProtoToLogEntry converts a TruncateRequest protobuf into a LogEntry
func TruncateProtoToLogEntry(pro *pb.TruncateRequest) LogEntry {
	return newLogEntry(typeTruncate, pro.Path, pro.Length)
}

// UnlinkProtoToLogEntry converts a UnlinkRequest protobuf into a LogEntry
func UnlinkProtoToLogEntry(pro *pb.UnlinkRequest) LogEntry {
	return newLogEntry(typeUnlink, pro.Path)
}

// UtimesProtoToLogEntry converts a UtimesRequest protobuf into a LogEntry
func UtimesProtoToLogEntry(pro *pb.UtimesRequest) LogEntry {
	return newLogEntry(typeUtimes, pro.Path, pro.AccessSeconds, pro.AccessNanoseconds, pro.ModifySeconds, pro.ModifyNanoseconds)
}

// WriteProtoToLogEntry converts a WriteRequest protobuf into a LogEntry
func WriteProtoToLogEntry(pro *pb.WriteRequest) LogEntry {
	return newLogEntry(typeWrite, pro.Path, pro.Data, pro.Offset, pro.Length)
}

// failedConversionCheck is a helper function for error checking for getEntryData
func failedConversionCheck(success bool, logType string, params ...interface{}) {
	if !success {
		log.Fatalln("Activity logger: Bad parameters for", logType, "LogEntry\n", params)
	}
}
