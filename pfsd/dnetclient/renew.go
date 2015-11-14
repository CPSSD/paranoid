package dnetclient

import (
	"errors"
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/discoverynetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

// Renew function. Will create a goroutine which will send renew to server
// 1/10 before expriration
func Renew() error {
	conn, err := grpc.Dial(globals.DiscoveryAddr, grpc.WithInsecure())
	if err != nil {
		log.Println("[D] [E] failed to dial discovery server at ", globals.DiscoveryAddr)
		return errors.New("[D] [E] failed to dial discovery server")
	}
	defer conn.Close()

	dclient := pb.NewDiscoveryNetworkClient(conn)
	pbNode := pb.Node{Ip: ThisNode.IP, Port: ThisNode.Port}

	_, err = dclient.Renew(context.Background(), &pb.JoinRequest{Node: &pbNode})
	if err != nil {
		log.Println("[D] [E] could not join")
		return errors.New("[D] [E] could not renew")
	}

	log.Println("[D] [I] Renewed discovery membership")
	return nil
}
