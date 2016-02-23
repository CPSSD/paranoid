//This file stores functions used to interface in or out of raft
package raft

import (
	"errors"
	pb "github.com/cpssd/paranoid/proto/raft"
	"google.golang.org/grpc"
	"net"
	"time"
)

//Starts a raft server given a listener, node information a directory to store information and a list of peers
func StartRaft(lis *net.Listener, nodeDetails Node, raftInfoDirectory string, peers []Node) (*RaftNetworkServer, *grpc.Server) {
	var opts []grpc.ServerOption
	srv := grpc.NewServer(opts...)
	raftServer := newRaftNetworkServer(nodeDetails, raftInfoDirectory, peers)
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
func (s *RaftNetworkServer) RequestAddLogEntry(entry *pb.Entry) error {
	s.addEntryLock.Lock()
	defer s.addEntryLock.Unlock()
	currentState := s.State.GetCurrentState()

	s.State.SetWaitingForApplied(true)
	defer s.State.SetWaitingForApplied(false)

	//Add entry to leaders Log
	if currentState == LEADER {
		err := s.addLogEntryLeader(entry)
		if err != nil {
			return err
		}
	} else if currentState == FOLLOWER {
		if s.State.GetLeaderId() != "" {
			err := s.sendLeaderLogEntry(entry)
			if err != nil {
				return err
			}
		} else {
			select {
			case <-time.After(20 * time.Second):
				return errors.New("Could not find a leader")
			case <-s.State.LeaderElected:
				if s.State.GetCurrentState() == LEADER {
					err := s.addLogEntryLeader(entry)
					if err != nil {
						return err
					}
				} else {
					err := s.sendLeaderLogEntry(entry)
					if err != nil {
						return err
					}
				}
			}
		}
	} else {
		count := 0
		for {
			count++
			if count > 40 {
				return errors.New("Could not find a leader")
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
				return err
			}
		} else {
			err := s.sendLeaderLogEntry(entry)
			if err != nil {
				return err
			}
		}
	}

	//Wait for the Log entry to be applied
	timer := time.NewTimer(ENTRY_APPLIED_TIMEOUT)
	for {
		select {
		case <-timer.C:
			return errors.New("Waited too long to commit Log entry")
		case entryIndex := <-s.State.EntryApplied:
			LogEntry, err := s.State.Log.GetLogEntry(entryIndex)
			if err != nil {
				Log.Fatal("Unable to get log entry:", err)
			}
			if LogEntry.Entry.Uuid == entry.Uuid {
				return nil
			}
		}
	}
	return nil
}

func (s *RaftNetworkServer) RequestChangeConfiguration(nodes []Node) error {
	Log.Info("Configuration change requested")
	entry := &pb.Entry{
		Type:    pb.Entry_ConfigurationChange,
		Uuid:    generateNewUUID(),
		Command: nil,
		Config: &pb.Configuration{
			Type:  pb.Configuration_FutureConfiguration,
			Nodes: convertNodesToProto(nodes),
		},
	}
	return s.RequestAddLogEntry(entry)
}
