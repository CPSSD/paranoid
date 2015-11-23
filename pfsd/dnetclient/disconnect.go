package dnetclient

import (
	"errors"
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

// Disconnect function used to disconnect from the server
func Disconnect() error {
	globals.Disconnecting = false
	conn, err := grpc.Dial(globals.DiscoveryAddr, grpc.WithInsecure())
	if err != nil {
		log.Println("ERROR: failed to dial discovery server at ", globals.DiscoveryAddr)
		return errors.New("Failed to dial discovery server")
	}
	defer conn.Close()

	dclient := pb.NewDiscoveryNetworkClient(conn)

	_, err = dclient.Disconnect(context.Background(),
		&pb.DisconnectRequest{Node: &pb.Node{Ip: ThisNode.IP, Port: ThisNode.Port}})
	if err != nil {
		log.Println("ERROR: could not send disconnect message")
		return errors.New("Could not send disconnect message")
	}

	return nil
}
