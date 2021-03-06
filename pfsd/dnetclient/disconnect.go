package dnetclient

import (
	"errors"
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
)

// Disconnect function used to disconnect from the server
func Disconnect(pool, password string) error {
	conn, err := dialDiscovery()
	if err != nil {
		Log.Error("Failed to dial discovery server at ", globals.DiscoveryAddr)
		return errors.New("failed to dial discovery server")
	}
	defer conn.Close()

	dclient := pb.NewDiscoveryNetworkClient(conn)

	_, err = dclient.Disconnect(context.Background(),
		&pb.DisconnectRequest{
			Pool:     pool,
			Password: password,
			Node: &pb.Node{
				Ip:   globals.ThisNode.IP,
				Port: globals.ThisNode.Port,
				Uuid: globals.ThisNode.UUID},
		})
	if err != nil {
		Log.Error("Could not send disconnect message")
		return errors.New("could not send disconnect message")
	}

	return nil
}
