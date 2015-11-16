package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func unlink(ips []globals.Node, path string) {
	for _, ipAddress := range ips {
		sendUnlinkRequest(ipAddress, path)
	}
}

func sendUnlinkRequest(ipAddress globals.Node, path string) {
	conn, err := grpc.Dial(ipAddress.IP+":"+ipAddress.Port, grpc.WithInsecure())
	if err != nil {
		log.Println("fail to dial: ", err)
	}

	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)

	response, err := client.Unlink(context.Background(), &pb.UnlinkRequest{path})
	if err != nil {
		log.Println("Unlink Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
	}
	log.Println(response)
}
