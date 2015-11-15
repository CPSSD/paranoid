package network

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func rename(ips []globals.Node, oldPath, newPath string) {
	for _, ipAddress := range ips {
		sendRenameRequest(ipAddress, oldPath, newPath)
	}
}

func sendRenameRequest(ipAddress globals.Node, oldPath, newPath string) {
	conn, err := grpc.Dial(ipAddress.IP+":"+ipAddress.Port, grpc.WithInsecure())
	if err != nil {
		log.Println("fail to dial: ", err)
	}

	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)

	response, err := client.Rename(context.Background(), &pb.RenameRequest{oldPath, newPath})
	if err != nil {
		log.Println("Rename Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
	}
	log.Println(response)
}
