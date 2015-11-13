package dnetclient

import (
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func Join(pool string) {
	conn, err := grpc.Dial(DiscoveryAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("[D] [E] failed to dial discovery server at ", DiscoveryAddr)
		return
	}
	defer conn.Close()

	dclient := pb.NewDiscoveryNetworkClient(conn)

	response, err := dclient.Join(context.Background(), &pb.JoinRequest{Pool: pool, Node: &thisNode})
	if err != nil {
		log.Print("[D] [E] could not join")
		return
	}

	resetInterval = response.ResetInterval
	Nodes = response.Nodes

}
