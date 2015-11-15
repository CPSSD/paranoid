package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func link(ips []globals.Node, oldPath, newPath string) {
	for _, ipAddress := range ips {
		sendLinkRequest(ipAddress, oldPath, newPath)
	}
}

func sendLinkRequest(ipAddress globals.Node, oldPath, newPath string) {
	conn, err := grpc.Dial(ipAddress.IP+":"+ipAddress.Port, grpc.WithInsecure())
	if err != nil {
		log.Println("fail to dial: ", err)
	}

	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)

	response, err := client.Link(context.Background(), &pb.LinkRequest{oldPath, newPath})
	if err != nil {
		log.Println("Link Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
	}
	log.Println(response)
}
