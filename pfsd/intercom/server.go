package intercom

import (
	"fmt"
	"github.com/cpssd/paranoid/pfsd/globals"
	"net"
	"net/rpc"
	"os"
	"path"
)

var lis net.Listener

type IntercomServer struct{}
type EmptyMessage struct{}

// Literally just a method for paranoid-cli to ping PFSD
func (s *IntercomServer) ConfirmUp(req *EmptyMessage, resp *EmptyMessage) error {
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
}

func ShutdownServer() error {
	err := lis.Close()
	if err != nil {
		return fmt.Errorf("failed to close listener: %s", err)
	}
	return nil
}
