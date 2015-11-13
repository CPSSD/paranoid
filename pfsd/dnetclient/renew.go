package dnetclient

import (
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"time"
)

// Renew function. Will create a goroutine which will send renew to server
// 1/10 before expriration
func Renew() {
	conn, err := grpc.Dial(DiscoveryAddr, grpc.WithInsecure())
	if err != nil {
		log.Println("[D] [E] failed to dial discovery server at ", DiscoveryAddr)
		return
	}
	defer conn.Close()

	dclient := pb.NewDiscoveryNetworkClient(conn)
	pbNode := pb.Node{ThisNode.IP, ThisNode.Port}

	go callRenew(dclient, pbNode)
}

func callRenew(dclient pb.DiscoveryNetworkClient, pbNode pb.Node) {
	for {
		_, err := dclient.Renew(context.Background(), &pb.JoinRequest{Pool: "_", Node: &pbNode})
		if err != nil {
			log.Println("[D] [E] could not join")
		}

		time.Sleep(resetInterval)
	}

}
