package raft

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	pb "github.com/cpssd/paranoid/proto/raft"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

const (
	SnapshotDirectory        string = "snapshots"
	CurrentSnapshotDirectory string = "nextsnapshot"
	NextShapshotDirectory    string = "nextsnapshot"
	SnapshotMetaFileName     string = "snapshotmeta"
	TarFileName              string = "snapshot.tar"
)

// Called every time raft network server starts up
// Creates the snapshot directory if it does not exist and removes any snapshot we were building that was not completed
func (s *RaftNetworkServer) setupSnapshotDirectory() {
	_, err := os.Stat(path.Join(s.raftInfoDirectory, SnapshotDirectory))
	if os.IsNotExist(err) {
		err := os.Mkdir(path.Join(s.raftInfoDirectory, SnapshotDirectory), 0700)
		if err != nil {
			Log.Fatal("failed to create snapshot directory:", err)
		}
	} else if err != nil {
		Log.Fatal("error acessing snapshot directory:", err)
	}

	_, err = os.Stat(path.Join(s.raftInfoDirectory, SnapshotDirectory, NextShapshotDirectory))
	if err == nil {
		err := os.RemoveAll(path.Join(s.raftInfoDirectory, SnapshotDirectory, NextShapshotDirectory))
		if err != nil {
			Log.Fatal("failed to discard old in progress snapshot")
		}
	}
}

type SnapShotInfo struct {
	LastIncludedIndex uint64 `json:"lastincludedindex"`
	LastIncludedTerm  uint64 `json:"lastincludedterm"`
}

func getSnapshotMetaInformation(snapShotPath string) (*SnapShotInfo, error) {
	metaFileContents, err := ioutil.ReadFile(path.Join(snapShotPath, SnapshotMetaFileName))
	snapShotInfo := &SnapShotInfo{}
	if err != nil {
		return snapShotInfo, fmt.Errorf("error reading raft meta information: %s", err)
	}

	err = json.Unmarshal(metaFileContents, &snapShotInfo)
	if err != nil {
		return snapShotInfo, fmt.Errorf("error reading raft meta information: %s", err)
	}
	return snapShotInfo, nil
}

func saveSnapshotMetaInformation(snapShotPath string, lastIncludedIndex, lastIncludedTerm uint64) error {
	snapShotInfo := &SnapShotInfo{
		LastIncludedIndex: lastIncludedIndex,
		LastIncludedTerm:  lastIncludedTerm,
	}

	snapShotInfoJson, err := json.Marshal(snapShotInfo)
	if err != nil {
		return fmt.Errorf("error saving snapshot meta information: %s", err)
	}

	err = ioutil.WriteFile(path.Join(snapShotPath, SnapshotMetaFileName), snapShotInfoJson, 0600)
	if err != nil {
		return fmt.Errorf("error saving snapshot meta information: %s", err)
	}

	return nil
}

func unpackTarFile(tarFilePath, directory string) error {
	untar := exec.Command("tar", "-xf", tarFilePath, "--directory="+directory)
	err := untar.Run()
	if err != nil {
		return fmt.Errorf("error unarchiving %s: %s", tarFilePath, err)
	}

	err = os.Rename(path.Join(directory, "contents-tar"), path.Join(directory, "contents"))
	if err != nil {
		return fmt.Errorf("error unpacking %s: %s", tarFilePath, err)
	}

	err = os.Rename(path.Join(directory, "names-tar"), path.Join(directory, "names"))
	if err != nil {
		return fmt.Errorf("error unpacking %s: %s", tarFilePath, err)
	}

	err = os.Rename(path.Join(directory, "inodes-tar"), path.Join(directory, "inodes"))
	if err != nil {
		return fmt.Errorf("error unpacking %s: %s", tarFilePath, err)
	}
	return nil
}

func tarSnapshot(snapshotDirectory string) error {
	contentsPath := path.Join(snapshotDirectory, "contents-tar")
	namesPath := path.Join(snapshotDirectory, "names-tar")
	inodesPath := path.Join(snapshotDirectory, "inodes-tar")
	configPath := path.Join(snapshotDirectory, PersistentConfigurationFileName)
	tarFilePath := path.Join(snapshotDirectory, TarFileName)

	err := os.Rename(path.Join(snapshotDirectory, "contents"), contentsPath)
	if err != nil {
		return fmt.Errorf("error creating tar file: %s", err)
	}

	err = os.Rename(path.Join(snapshotDirectory, "inodes"), inodesPath)
	if err != nil {
		return fmt.Errorf("error creating tar file: %s", err)
	}

	err = os.Rename(path.Join(snapshotDirectory, "names"), namesPath)
	if err != nil {
		return fmt.Errorf("error creating tar file: %s", err)
	}

	tar := exec.Command("tar", "cf", tarFilePath, contentsPath, namesPath, inodesPath, configPath)
	err = tar.Run()
	if err != nil {
		return fmt.Errorf("error creating tar file: %s", err)
	}

	return nil
}

func startNextSnapshotWithCurrent(currentSnapshot, nextSnapshot string) error {
	err := unpackTarFile(path.Join(currentSnapshot, TarFileName), nextSnapshot)
	if err != nil {
		return fmt.Errorf("error starting new snapshot from current snapshot:", err)
	}

	return nil
}

func copyFile(originalPath, copyPath string) error {
	contents, err := ioutil.ReadFile(originalPath)
	if err != nil {
		return fmt.Errorf("error copying file: %s", err)
	}

	err = ioutil.WriteFile(copyPath, contents, 0600)
	if err != nil {
		return fmt.Errorf("error copying file: %s", err)
	}

	return nil
}

func (s *RaftNetworkServer) startNextSnapshot(nextSnapshot string) error {
	err := os.Mkdir(path.Join(nextSnapshot, "contents"), 0700)
	if err != nil {
		return fmt.Errorf("error starting next snapshot: %s", err)
	}

	err = os.Mkdir(path.Join(nextSnapshot, "inodes"), 0700)
	if err != nil {
		return fmt.Errorf("error starting next snapshot: %s", err)
	}

	err = os.Mkdir(path.Join(nextSnapshot, "names"), 0700)
	if err != nil {
		return fmt.Errorf("error starting next snapshot: %s", err)
	}

	err = copyFile(path.Join(s.raftInfoDirectory, OriginalConfigurationFileName), path.Join(nextSnapshot, PersistentConfigurationFileName))
	if err != nil {
		return fmt.Errorf("error starting next snapshot: %s", err)
	}

	return nil
}

func (s *RaftNetworkServer) applyLogUpdates(snapshotDirectory string, startIndex, endIndex uint64) (lastIncludedTerm uint64, err error) {
	snapshotConfig := newConfiguration(snapshotDirectory, nil, s.nodeDetails, false)

	if startIndex > endIndex {
		return 0, errors.New("no log entries to apply")
	}

	for i := startIndex; i <= endIndex; i++ {
		logEntry, err := s.State.Log.GetLogEntry(i)
		if err != nil {
			return 0, fmt.Errorf("unable to apply log entry: %s", err)
		}
		if i == endIndex {
			lastIncludedTerm = logEntry.Term
		}

		if logEntry.Entry.Type == pb.Entry_StateMachineCommand {
			libpfsCommand := logEntry.Entry.GetCommand()
			if libpfsCommand == nil {
				return 0, errors.New("unable to apply log entry with empty command field")
			}

			result := PerformLibPfsCommand(snapshotDirectory, libpfsCommand)
			if result.Code == returncodes.EUNEXPECTED {
				return 0, fmt.Errorf("error applying log entry: %s", result.Err)
			}
		} else if logEntry.Entry.Type == pb.Entry_ConfigurationChange {
			config := logEntry.Entry.GetConfig()
			if config == nil {
				return 0, errors.New("unable to apply log entry with empty config field")
			}

			if config.Type == pb.Configuration_CurrentConfiguration {
				snapshotConfig.UpdateCurrentConfiguration(protoNodesToNodes(config.Nodes), 0)
			} else {
				snapshotConfig.NewFutureConfiguration(protoNodesToNodes(config.Nodes), 0)
			}
		} else {
			return 0, fmt.Errorf("unable to snapshot command type %s", logEntry.Entry.Type)
		}
	}

	return lastIncludedTerm, nil
}

func (s *RaftNetworkServer) createSnapshot(lastIncludedIndex uint64) error {
	defer s.Wait.Done()

	currentSnapshot := path.Join(s.raftInfoDirectory, SnapshotDirectory, CurrentSnapshotDirectory)
	nextSnapshot := path.Join(s.raftInfoDirectory, SnapshotDirectory, NextShapshotDirectory)

	if s.State.GetPerformingSnapshot() == true {
		return errors.New("snapshot creation already in progress")
	}
	s.State.SetPerformingSnapshot(true)

	_, err := os.Stat(nextSnapshot)
	if !os.IsNotExist(err) {
		return errors.New("snapshot creation already in progress")
	}

	startLogIndex := uint64(0)

	_, err = os.Stat(currentSnapshot)
	if err == nil {
		metaInfo, err := getSnapshotMetaInformation(currentSnapshot)
		if err != nil {
			return err
		}
		startLogIndex = metaInfo.LastIncludedIndex + 1

		err = startNextSnapshotWithCurrent(currentSnapshot, nextSnapshot)
		if err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("could not access current snapshot directory: %s", err)
	} else {
		err = s.startNextSnapshot(nextSnapshot)
		if err != nil {
			return err
		}
	}

	lastIncludedTerm, err := s.applyLogUpdates(nextSnapshot, startLogIndex, lastIncludedIndex)
	if err != nil {
		return err
	}

	err = saveSnapshotMetaInformation(nextSnapshot, lastIncludedIndex, lastIncludedTerm)
	if err != nil {
		return err
	}

	err = tarSnapshot(nextSnapshot)
	if err != nil {
		return err
	}

	err = os.Rename(nextSnapshot, currentSnapshot)
	if err != nil {
		return err
	}

	return nil
}

func (s *RaftNetworkServer) InstallSnapshot(ctx context.Context, req *pb.SnapshotRequest) (*pb.SnapshotResponse, error) {
	if req.Term < s.State.GetCurrentTerm() {
		return &pb.SnapshotResponse{s.State.GetCurrentTerm()}, nil
	}

	return &pb.SnapshotResponse{s.State.GetCurrentTerm()}, nil
}
