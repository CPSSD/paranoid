package intercom

import (
	"fmt"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/pnetserver"
	"net"
	"net/rpc"
	"os"
	"path"
	"strconv"
	"time"
)

var lis net.Listener

type PFSDStatus int

const (
	FOLLOWER PFSDStatus = iota
	CANDIDATE
	LEADER
	RAFT_INACTIVE
)

type IntercomServer struct{}
type EmptyMessage struct{}

type StatusResponse struct {
	Uptime    time.Duration
	Status    PFSDStatus
	TLSActive bool
	Port      int
}

// Literally just a method for paranoid-cli to ping PFSD
func (s *IntercomServer) ConfirmUp(req *EmptyMessage, resp *EmptyMessage) error {
	return nil
}

// Provides health data for the current node.
func (s *IntercomServer) Status(req *EmptyMessage, resp *StatusResponse) error {
	thisport, err := strconv.Atoi(globals.ThisNode.Port)
	if err != nil {
		Log.Error("Could not convert globals.ThisNode.Port to int.")
		return fmt.Errorf("failed converting globals.ThisNode.Port to int: %s", err)
	}

	resp.Uptime = time.Since(globals.BootTime)
	resp.Status = PFSDStatus(pnetserver.RaftNetworkServer.State.GetCurrentState())
	resp.TLSActive = globals.TLSEnabled
	resp.Port = thisport
	return nil
}

func RunServer(metaDir string) {
	socketPath := path.Join(metaDir, "intercom.sock")
	err := os.Remove(socketPath)
	if err != nil && !os.IsNotExist(err) {
		Log.Fatalf("Failed to remove %s: %s\n", socketPath, err)
	}
	server := new(IntercomServer)
	rpc.Register(server)
	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		Log.Fatalf("Failed to listen on %s: %s\n", socketPath, err)
	}
	globals.Wait.Add(1)
	go func() {
		defer globals.Wait.Done()
		Log.Info("Internal communication server listening on", socketPath)
		rpc.Accept(lis)
	}()
	globals.BootTime = time.Now()
	Log.Info(globals.BootTime)
}

func ShutdownServer() error {
	err := lis.Close()
	if err != nil {
		return fmt.Errorf("failed to close listener: %s", err)
	}
	return nil
}
