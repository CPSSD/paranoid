package dnetclient

import (
	"errors"
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

// Join function to call in order to join the server
func Join(pool string) error {
	conn, err := dialDiscovery()
	if err != nil {
		return errors.New("Failed to dial discovery server")
	}
	defer conn.Close()

	dclient := pb.NewDiscoveryNetworkClient(conn)

	response, err := dclient.Join(context.Background(),
		&pb.JoinRequest{
			Pool: pool,
			Node: &pb.Node{
				Ip:         ThisNode.IP,
				Port:       ThisNode.Port,
				CommonName: ThisNode.CommonName,
				Uuid:       ThisNode.UUID,
			},
		})
	if err != nil {
		return errors.New("Could not join the pool")
	}

	interval := response.ResetInterval / 10 * 9
	globals.ResetInterval, err = time.ParseDuration(strconv.FormatInt(interval, 10) + "ms")
	if err != nil {
		Log.Error("Invalid renew interval.", err)
	}

	peerList := "Currently Connected: "
	for _, node := range response.Nodes {
		peerList += node.Ip + ":" + node.Port + ", "
		globals.Nodes.Add(globals.Node{IP: node.Ip, Port: node.Port})
	}
	Log.Info(peerList)

	Log.Info("Successfully joined discovery network")

	return nil
}
