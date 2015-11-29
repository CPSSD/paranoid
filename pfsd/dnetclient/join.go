package dnetclient

import (
	"errors"
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"time"
)

// Join function to call in order to join the server
func Join(pool string) error {
	conn, err := dialDiscovery()
	if err != nil {
		log.Println("ERROR: failed to dial discovery server at ", globals.DiscoveryAddr)
		return errors.New("Failed to dial discovery server")
	}
	defer conn.Close()

	dclient := pb.NewDiscoveryNetworkClient(conn)

	response, err := dclient.Join(context.Background(),
		&pb.JoinRequest{Pool: pool, Node: &pb.Node{Ip: ThisNode.IP, Port: ThisNode.Port}})
	if err != nil {
		log.Println("ERROR: could not join discovery server")
		return errors.New("Could not join the pool")
	}

	interval := response.ResetInterval / 10 * 9
	globals.ResetInterval, _ = time.ParseDuration(strconv.FormatInt(interval, 10) + "ms")

	for _, node := range response.Nodes {
		globals.Nodes.Add(globals.Node{IP: node.Ip, Port: node.Port})
	}

	log.Println("INFO: Successfully joined")
	return nil
}
