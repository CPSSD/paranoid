//This file stores functions used to interface in or out of raft
package raft

import (
	"errors"
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	pb "github.com/cpssd/paranoid/proto/raft"
	"google.golang.org/grpc"
	"net"
	"time"
)

const (
	TYPE_WRITE uint32 = iota
)

type StateMachineResult struct {
	Code         int
	Err          error
	BytesWritten int
}

type EntryAppliedInfo struct {
	Index  uint64
	Result *StateMachineResult
}

//Starts a raft server given a listener, node information a directory to store information
//A start configuration can be given for testing or for the first node in a cluster
func StartRaft(lis *net.Listener, nodeDetails Node, pfsDirectory, raftInfoDirectory string,
	startConfiguration *StartConfiguration) (*RaftNetworkServer, *grpc.Server) {

	var opts []grpc.ServerOption
	srv := grpc.NewServer(opts...)
	raftServer := newRaftNetworkServer(nodeDetails, pfsDirectory, raftInfoDirectory, startConfiguration)
	pb.RegisterRaftNetworkServer(srv, raftServer)
	raftServer.Wait.Add(1)
	go func() {
		Log.Info("RaftNetworkServer started")
		err := srv.Serve(*lis)
		if err != nil {
			Log.Error("Error running RaftNetworkServer", err)
		}
	}()
	return raftServer, srv
}

//A request to add a Log entry from a client. If the node is not the leader, it must forward the request to the leader.
//Only return once the request has been commited to the State machine
func (s *RaftNetworkServer) RequestAddLogEntry(entry *pb.Entry) (error, *StateMachineResult) {
	s.addEntryLock.Lock()
	defer s.addEntryLock.Unlock()
	currentState := s.State.GetCurrentState()

	s.State.SetWaitingForApplied(true)
	defer s.State.SetWaitingForApplied(false)

	//Add entry to leaders Log
	if currentState == LEADER {
		err := s.addLogEntryLeader(entry)
		if err != nil {
			return err, nil
		}
	} else if currentState == FOLLOWER {
		if s.State.GetLeaderId() != "" {
			err := s.sendLeaderLogEntry(entry)
			if err != nil {
				return err, nil
			}
		} else {
			select {
			case <-time.After(20 * time.Second):
				return errors.New("Could not find a leader"), nil
			case <-s.State.LeaderElected:
				if s.State.GetCurrentState() == LEADER {
					err := s.addLogEntryLeader(entry)
					if err != nil {
						return err, nil
					}
				} else {
					err := s.sendLeaderLogEntry(entry)
					if err != nil {
						return err, nil
					}
				}
			}
		}
	} else {
		count := 0
		for {
			count++
			if count > 40 {
				return errors.New("Could not find a leader"), nil
			}
			time.Sleep(500 * time.Millisecond)
			currentState = s.State.GetCurrentState()
			if currentState != CANDIDATE {
				break
			}
		}
		if currentState == LEADER {
			err := s.addLogEntryLeader(entry)
			if err != nil {
				return err, nil
			}
		} else {
			err := s.sendLeaderLogEntry(entry)
			if err != nil {
				return err, nil
			}
		}
	}

	//Wait for the Log entry to be applied
	timer := time.NewTimer(ENTRY_APPLIED_TIMEOUT)
	for {
		select {
		case <-timer.C:
			return errors.New("Waited too long to commit Log entry"), nil
		case appliedEntry := <-s.State.EntryApplied:
			LogEntry, err := s.State.Log.GetLogEntry(appliedEntry.Index)
			if err != nil {
				Log.Fatal("Unable to get log entry:", err)
			}
			if LogEntry.Entry.Uuid == entry.Uuid {
				return nil, appliedEntry.Result
			}
		}
	}
	return errors.New("Waited too long to commit Log entry"), nil
}

func (s *RaftNetworkServer) RequestWriteCommand(filePath string, offset, length uint64,
	data []byte) (returnCode int, returnError error, bytesWrote int) {
	entry := &pb.Entry{
		Type: pb.Entry_StateMachineCommand,
		Uuid: generateNewUUID(),
		Command: &pb.StateMachineCommand{
			Type:   TYPE_WRITE,
			Path:   filePath,
			Data:   data,
			Offset: offset,
			Length: length,
		},
	}
	err, stateMachineResult := s.RequestAddLogEntry(entry)
	if err != nil {
		return returncodes.EBUSY, err, 0
	}
	return stateMachineResult.Code, stateMachineResult.Err, stateMachineResult.BytesWritten
}

func (s *RaftNetworkServer) RequestChangeConfiguration(nodes []Node) error {
	Log.Info("Configuration change requested")
	entry := &pb.Entry{
		Type: pb.Entry_ConfigurationChange,
		Uuid: generateNewUUID(),
		Config: &pb.Configuration{
			Type:  pb.Configuration_FutureConfiguration,
			Nodes: convertNodesToProto(nodes),
		},
	}
	err, _ := s.RequestAddLogEntry(entry)
	return err
}

func performLibPfsCommand(directory string, command *pb.StateMachineCommand) *StateMachineResult {
	switch command.Type {
	case TYPE_WRITE:
		code, err, bytesWritten := commands.WriteCommand(directory, command.Path, int64(command.Offset), int64(command.Length), command.Data)
		return &StateMachineResult{code, err, bytesWritten}
	}
	Log.Fatal("Unrecognised command type")
	return nil
}
