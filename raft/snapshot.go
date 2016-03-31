package raft

import (
	"encoding/json"
	"errors"
	"fmt"
	libpfs "github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
	pb "github.com/cpssd/paranoid/proto/raft"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"
)

const (
	SnapshotDirectory        string = "snapshots"
	CurrentSnapshotDirectory string = "currentsnapshot"
	SnapshotMetaFileName     string = "snapshotmeta"
	TarFileName              string = "snapshot.tar"
)

const (
	SNAPSHOT_INTERVAL         time.Duration = 1 * time.Minute
	SNAPSHOT_LOGSIZE          uint64        = 2 * 1024 * 1024 //2 MegaBytes
	SNAPSHOT_CHUNK_SIZE       int64         = 1024
	MAX_INSTALLSNAPSHOT_FAILS int           = 10
)

// Called every time raft network server starts up
// Makes sure snapshot directory exists and we have access to it
func (s *RaftNetworkServer) setupSnapshotDirectory() {
	_, err := os.Stat(path.Join(s.raftInfoDirectory, SnapshotDirectory))
	if os.IsNotExist(err) {
		err := os.Mkdir(path.Join(s.raftInfoDirectory, SnapshotDirectory), 0700)
		if err != nil {
			Log.Fatal("failed to create snapshot directory:", err)
		}
	} else if err != nil {
		Log.Fatal("error accessing snapshot directory:", err)
	}
}

type SnapShotInfo struct {
	LastIncludedIndex uint64 `json:"lastincludedindex"`
	LastIncludedTerm  uint64 `json:"lastincludedterm"`
	SelfCreated       bool   `json:"selfcreated"`
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

func saveSnapshotMetaInformation(snapShotPath string, lastIncludedIndex, lastIncludedTerm uint64, selfCreated bool) error {
	snapShotInfo := &SnapShotInfo{
		LastIncludedIndex: lastIncludedIndex,
		LastIncludedTerm:  lastIncludedTerm,
		SelfCreated:       selfCreated,
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

	err = os.RemoveAll(path.Join(directory, "contents"))
	if err != nil {
		return fmt.Errorf("error unpacking %s: %s", tarFilePath, err)
	}
	err = os.Rename(path.Join(directory, "contents-tar"), path.Join(directory, "contents"))
	if err != nil {
		return fmt.Errorf("error unpacking %s: %s", tarFilePath, err)
	}

	err = os.RemoveAll(path.Join(directory, "names"))
	if err != nil {
		return fmt.Errorf("error unpacking %s: %s", tarFilePath, err)
	}
	err = os.Rename(path.Join(directory, "names-tar"), path.Join(directory, "names"))
	if err != nil {
		return fmt.Errorf("error unpacking %s: %s", tarFilePath, err)
	}

	err = os.RemoveAll(path.Join(directory, "inodes"))
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
	err := os.Rename(path.Join(snapshotDirectory, "contents"), path.Join(snapshotDirectory, "contents-tar"))
	if err != nil {
		return fmt.Errorf("error creating tar file: %s", err)
	}

	err = os.Rename(path.Join(snapshotDirectory, "inodes"), path.Join(snapshotDirectory, "inodes-tar"))
	if err != nil {
		return fmt.Errorf("error creating tar file: %s", err)
	}

	err = os.Rename(path.Join(snapshotDirectory, "names"), path.Join(snapshotDirectory, "names-tar"))
	if err != nil {
		return fmt.Errorf("error creating tar file: %s", err)
	}

	tar := exec.Command("tar", "--directory="+snapshotDirectory, "-cf", path.Join(snapshotDirectory, TarFileName),
		"contents-tar", "names-tar", "inodes-tar", PersistentConfigurationFileName)
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

	err = os.Mkdir(path.Join(nextSnapshot, "meta"), 0700)
	if err != nil {
		return fmt.Errorf("error starting next snapshot: %s", err)
	}

	_, err = os.Create(path.Join(nextSnapshot, "meta", "lock"))
	if err != nil {
		return fmt.Errorf("error starting next snapshot: %s", err)
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

	err = os.Mkdir(path.Join(nextSnapshot, "meta"), 0700)
	if err != nil {
		return fmt.Errorf("error starting next snapshot: %s", err)
	}

	_, err = os.Create(path.Join(nextSnapshot, "meta", "lock"))
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

//performCleanup is used to clean up tempory files used in snapshot creation
func performCleanup(snapshotPath string) error {
	err := os.RemoveAll(path.Join(snapshotPath, "contents-tar"))
	if err != nil {
		return fmt.Errorf("error cleaning up temporay files: %s", err)
	}

	err = os.RemoveAll(path.Join(snapshotPath, "inodes-tar"))
	if err != nil {
		return fmt.Errorf("error cleaning up temporay files: %s", err)
	}

	err = os.RemoveAll(path.Join(snapshotPath, "names-tar"))
	if err != nil {
		return fmt.Errorf("error cleaning up temporay files: %s", err)
	}

	err = os.RemoveAll(path.Join(snapshotPath, "meta"))
	if err != nil {
		return fmt.Errorf("error cleaning up temporay files: %s", err)
	}

	err = os.Remove(path.Join(snapshotPath, PersistentConfigurationFileName))
	if err != nil {
		return fmt.Errorf("error cleaning up temporay files: %s", err)
	}
	return nil
}

func (s *RaftNetworkServer) CreateSnapshot(lastIncludedIndex uint64) (err error) {
	currentSnapshot := path.Join(s.raftInfoDirectory, SnapshotDirectory, CurrentSnapshotDirectory)
	nextSnapshot := path.Join(s.raftInfoDirectory, SnapshotDirectory, generateNewUUID())

	if s.State.GetPerformingSnapshot() == true {
		return errors.New("snapshot creation already in progress")
	}
	s.State.SetPerformingSnapshot(true)
	defer s.State.SetPerformingSnapshot(false)

	err = os.Mkdir(nextSnapshot, 0700)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			cleanuperror := os.RemoveAll(nextSnapshot)
			if cleanuperror != nil {
				Log.Error("error removing temporary snapshot creation files:", cleanuperror)
			}
		}
	}()

	startLogIndex := uint64(1)

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

	err = tarSnapshot(nextSnapshot)
	if err != nil {
		return err
	}

	err = saveSnapshotMetaInformation(nextSnapshot, lastIncludedIndex, lastIncludedTerm, true)
	if err != nil {
		return err
	}

	err = performCleanup(nextSnapshot)
	if err != nil {
		return err
	}

	s.State.NewSnapshotCreated <- true

	return nil
}

//Revert the statemachine state to the snapshot state.
//Remove all log entries
func (s *RaftNetworkServer) RevertToSnapshot(snapshotPath string) error {
	s.State.ApplyEntryLock.Lock()
	defer s.State.ApplyEntryLock.Unlock()

	snapshotMeta, err := getSnapshotMetaInformation(snapshotPath)
	if err != nil {
		if err != nil {
			return fmt.Errorf("error reverting to snapshot: %s", err)
		}
	}

	err = libpfs.GetFileSystemLock(s.State.pfsDirectory, libpfs.ExclusiveLock)
	if err != nil {
		return fmt.Errorf("error reverting to snapshot: %s", err)
	}

	defer func() {
		err := libpfs.UnLockFileSystem(s.State.pfsDirectory)
		if err != nil {
			Log.Fatal("error reverting to snapshot: %s", err)
		}
	}()

	err = unpackTarFile(path.Join(snapshotPath, TarFileName), s.State.pfsDirectory)
	if err != nil {
		return fmt.Errorf("error reverting to snapshot: %s", err)
	}

	err = os.Rename(path.Join(s.State.pfsDirectory, PersistentConfigurationFileName), path.Join(s.raftInfoDirectory, PersistentConfigurationFileName))
	if err != nil {
		return fmt.Errorf("error reverting to snapshot: %s", err)
	}

	err = s.State.Configuration.UpdateFromConfigurationFile(path.Join(s.raftInfoDirectory, PersistentConfigurationFileName), snapshotMeta.LastIncludedIndex)
	if err != nil {
		return fmt.Errorf("error reverting to snapshot: %s", err)
	}

	currentSnapshot := path.Join(s.raftInfoDirectory, SnapshotDirectory, CurrentSnapshotDirectory)
	if currentSnapshot != snapshotPath {
		err = os.Rename(snapshotPath, currentSnapshot)
		if err != nil {
			return fmt.Errorf("error reverting to snapshot: %s", err)
		}
	}

	s.State.Log.DiscardAllLogEntries(snapshotMeta.LastIncludedIndex, snapshotMeta.LastIncludedTerm)
	s.State.SetLastApplied(snapshotMeta.LastIncludedIndex)
	s.State.SetCommitIndex(snapshotMeta.LastIncludedIndex)

	return nil
}

func (s *RaftNetworkServer) InstallSnapshot(ctx context.Context, req *pb.SnapshotRequest) (*pb.SnapshotResponse, error) {
	if req.Term < s.State.GetCurrentTerm() {
		return &pb.SnapshotResponse{s.State.GetCurrentTerm()}, nil
	}

	snapshotPath := path.Join(s.raftInfoDirectory, SnapshotDirectory, req.LeaderId+strconv.FormatUint(req.LastIncludedIndex, 10))
	if req.Offset == 0 {
		err := os.RemoveAll(snapshotPath)
		if err != nil {
			Log.Error("Error recieving snapshot:", err)
			return &pb.SnapshotResponse{s.State.GetCurrentTerm()}, err
		}

		err = os.Mkdir(snapshotPath, 0700)
		if err != nil {
			Log.Error("Error recieving snapshot:", err)
			return &pb.SnapshotResponse{s.State.GetCurrentTerm()}, err
		}

		snapshotFile, err := os.Create(path.Join(snapshotPath, TarFileName))
		if err != nil {
			Log.Error("Error recieving snapshot:", err)
			return &pb.SnapshotResponse{s.State.GetCurrentTerm()}, err
		}
		snapshotFile.Close()
	}

	snapshotFile, err := os.OpenFile(path.Join(snapshotPath, TarFileName), os.O_WRONLY, 0600)
	if err != nil {
		Log.Error("Error recieving snapshot:", err)
		return &pb.SnapshotResponse{s.State.GetCurrentTerm()}, err
	}
	defer snapshotFile.Close()

	bytesWritten, err := snapshotFile.WriteAt(req.Data, int64(req.Offset))
	if err != nil {
		Log.Error("Error recieving snapshot:", err)
		return &pb.SnapshotResponse{s.State.GetCurrentTerm()}, err
	}
	if bytesWritten != len(req.Data) {
		Log.Error("Error recieving snapshot: incorrect number of bytes written to snapshot file")
		return &pb.SnapshotResponse{s.State.GetCurrentTerm()}, errors.New("incorrect number of bytes written to snapshot file")
	}

	if req.Done == false {
		return &pb.SnapshotResponse{s.State.GetCurrentTerm()}, nil
	}

	saveSnapshotMetaInformation(snapshotPath, req.LastIncludedIndex, req.LastIncludedTerm, false)
	s.State.NewSnapshotCreated <- true

	return &pb.SnapshotResponse{s.State.GetCurrentTerm()}, nil
}

func (s *RaftNetworkServer) sendSnapshot(node *Node) {
	defer s.Wait.Done()
	defer s.State.DecrementSnapshotCounter()
	defer s.State.Configuration.SetSendingSnapshot(node.NodeID, false)

	conn, err := s.Dial(node, HEARTBEAT_TIMEOUT)
	if err != nil {
		Log.Error("error sending snapshot:", err)
		return
	}
	defer conn.Close()

	client := pb.NewRaftNetworkClient(conn)
	currentSnapshot := path.Join(s.raftInfoDirectory, SnapshotDirectory, CurrentSnapshotDirectory)
	snapshotMeta, err := getSnapshotMetaInformation(currentSnapshot)
	if err != nil {
		Log.Error("Error sending snapshot:", err)
		return
	}

	snapshotFile, err := os.Open(path.Join(currentSnapshot, TarFileName))
	if err != nil {
		Log.Error("Error sending snapshot:", err)
		return
	}
	defer snapshotFile.Close()

	snapshotChunk := make([]byte, SNAPSHOT_CHUNK_SIZE)
	snapshotFileOffset := int64(0)
	installRequestsFailed := 0

	for {
		select {
		case _, ok := <-s.Quit:
			if !ok {
				s.QuitChannelClosed = true
				Log.Info("Stop sending snapshot to:", node.String())
				return
			}
		default:
			if s.State.GetCurrentState() != LEADER {
				Log.Info("Ceasing sending snapshot due to state change")
				return
			}

			done := false
			bytesRead, err := snapshotFile.ReadAt(snapshotChunk, snapshotFileOffset)
			if err != nil {
				if err == io.EOF {
					done = true
				} else {
					Log.Error("Error sending snapshot:", err)
					return
				}
			}

			response, err := client.InstallSnapshot(context.Background(), &pb.SnapshotRequest{
				Term:              s.State.GetCurrentTerm(),
				LeaderId:          s.nodeDetails.NodeID,
				LastIncludedIndex: snapshotMeta.LastIncludedIndex,
				LastIncludedTerm:  snapshotMeta.LastIncludedTerm,
				Offset:            uint64(snapshotFileOffset),
				Data:              snapshotChunk[:bytesRead],
				Done:              done,
			})
			if err == nil {
				if response.Term > s.State.GetCurrentTerm() {
					s.State.StopLeading <- true
					return
				}
				if done {
					Log.Info("Sucessfully send complete snapshot to:", node.String())
					s.State.Configuration.SetNextIndex(node.NodeID, snapshotMeta.LastIncludedIndex+1)
					return
				}
				snapshotFileOffset = snapshotFileOffset + int64(bytesRead)
			} else {
				if installRequestsFailed > MAX_INSTALLSNAPSHOT_FAILS {
					Log.Error("InstallSnapshot request failed repeatedly:", err)
					return
				} else {
					Log.Warn("InstallSnapshot request failed:", err)
					installRequestsFailed++
				}
			}
		}
	}
}

//Update the current snapshot to the most recent snapshot available and remove all incomplete snapshots
func (s *RaftNetworkServer) updateCurrentSnapshot() error {
	snapshots, err := ioutil.ReadDir(path.Join(s.raftInfoDirectory, SnapshotDirectory))
	if err != nil {
		return fmt.Errorf("unable to update current snapshot: %s", err)
	}

	currentSnapshot := path.Join(s.raftInfoDirectory, SnapshotDirectory, CurrentSnapshotDirectory)
	currentSnapshotMeta, err := getSnapshotMetaInformation(currentSnapshot)
	if err != nil {
		currentSnapshotMeta = &SnapShotInfo{
			LastIncludedIndex: 0,
		}
	}

	mostRecentSnapshot := ""
	mostRecentSnapshotMeta := currentSnapshotMeta

	for i := 0; i < len(snapshots); i++ {
		snapshotPath := path.Join(s.raftInfoDirectory, SnapshotDirectory, snapshots[i].Name())
		snapshotMeta, err := getSnapshotMetaInformation(snapshotPath)
		if err != nil {
			Log.Warn("error updating current snapshot:", err)
		} else {
			if snapshotMeta.LastIncludedIndex > mostRecentSnapshotMeta.LastIncludedIndex {
				mostRecentSnapshot = snapshotPath
				mostRecentSnapshotMeta = snapshotMeta
			}
		}
	}

	if mostRecentSnapshot != "" {
		err = os.RemoveAll(currentSnapshot)
		if err != nil {
			return fmt.Errorf("unable to update current snapshot: %s", err)
		}

		err = os.Rename(mostRecentSnapshot, currentSnapshot)
		if err != nil {
			Log.Fatal("Failed to rename new snapshot after deleteing current snapshot:", err)
		}

		if mostRecentSnapshotMeta.SelfCreated {
			s.State.Log.DiscardLogEntriesBefore(mostRecentSnapshotMeta.LastIncludedIndex, mostRecentSnapshotMeta.LastIncludedTerm)
		} else {
			err = s.RevertToSnapshot(currentSnapshot)
			if err != nil {
				Log.Fatal("Update current snapshot failed:", err)
			}
		}

		for i := 0; i < len(snapshots); i++ {
			snapshotPath := path.Join(s.raftInfoDirectory, SnapshotDirectory, snapshots[i].Name())
			if snapshotPath != currentSnapshot && snapshotPath != mostRecentSnapshot {
				err = os.RemoveAll(snapshotPath)
				if err != nil {
					Log.Warn("error updating current snapshot:", err)
				}
			}
		}
	}

	return nil
}

func (s *RaftNetworkServer) manageSnapshoting() {
	defer s.Wait.Done()
	snapshotTimer := time.NewTimer(SNAPSHOT_INTERVAL)
	for {
		select {
		case _, ok := <-s.Quit:
			if !ok {
				s.QuitChannelClosed = true
				Log.Info("Exiting snapshot managment loop")
				return
			}
		case <-snapshotTimer.C:
			if s.State.GetPerformingSnapshot() == false {
				if s.State.Log.GetLogSizeBytes() > SNAPSHOT_LOGSIZE {
					s.Wait.Add(1)
					go func() {
						defer s.Wait.Done()
						err := s.CreateSnapshot(s.State.GetLastApplied())
						if err != nil {
							Log.Error("manage snapshotting:", err)
						}
					}()
				}
			}
			snapshotTimer.Reset(SNAPSHOT_INTERVAL)
		case <-s.State.SnapshotCounterAtZero:
			err := s.updateCurrentSnapshot()
			if err != nil {
				Log.Error("manage snapshoting:", err)
			}
		case <-s.State.NewSnapshotCreated:
			Log.Info("New Snapshot Created")
			if s.State.GetSnapshotCounterValue() == 0 {
				err := s.updateCurrentSnapshot()
				if err != nil {
					Log.Error("manage snapshoting: %s", err)
				}
			}
		case node := <-s.State.SendSnapshot:
			if s.State.Configuration.GetSendingSnapshot(node.NodeID) == false {
				Log.Info("Send current snapshot to", node)
				s.State.Configuration.SetSendingSnapshot(node.NodeID, true)
				s.Wait.Add(1)
				s.State.IncrementSnapshotCounter()
				go s.sendSnapshot(&node)
			}
		}
	}
}
