package dnetclient

import (
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"time"
)

// Join function to call in order to join the server
func Join(pool string) {
	conn, err := grpc.Dial(DiscoveryAddr, grpc.WithInsecure())
	if err != nil {
		log.Println("[D] [E] failed to dial discovery server at ", DiscoveryAddr)
		return
	}
	defer conn.Close()

	dclient := pb.NewDiscoveryNetworkClient(conn)

	response, err := dclient.Join(context.Background(),
		&pb.JoinRequest{Pool: pool, Node: &pb.Node{ThisNode.IP, ThisNode.Port}})
	if err != nil {
		log.Println("[D] [E] could not join")
		return
	}

	interval := response.ResetInterval / 10 * 9
	resetInterval, _ = time.ParseDuration(strconv.FormatInt(interval, 10) + "ms")

	for _, node := range response.Nodes {
		Nodes = append(Nodes, Node{node.Ip, node.Port})
	}
}
