package dnetclient

import (
	"errors"
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
)

// Renew function. Will create a goroutine which will send renew to server
// 1/10 before expriration
func Renew() error {
	conn, err := dialDiscovery()
	if err != nil {
		Log.Error("Failed to dial discovery server at", globals.DiscoveryAddr)
		return errors.New("failed to dial discovery server")
	}
	defer conn.Close()

	dclient := pb.NewDiscoveryNetworkClient(conn)
	pbNode := pb.Node{Ip: ThisNode.IP, Port: ThisNode.Port, Uuid: ThisNode.UUID}

	_, err = dclient.Renew(context.Background(), &pb.JoinRequest{Node: &pbNode})
	if err != nil {
		return errors.New("could not renew discovery membership")
	}

	Log.Info("Renewed discovery membership")
	return nil
}
