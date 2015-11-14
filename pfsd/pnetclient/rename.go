package network

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func rename(ips []globals.Node, oldPath, newPath string) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	for _, ipAddress := range ips {
		sendRenameRequest(ipAddress, oldPath, newPath, opts)
	}
}

func sendRenameRequest(ipAddress globals.Node, oldPath, newPath string, opts []grpc.DialOption) {
	conn, err := grpc.Dial(ipAddress.IP+":"+ipAddress.Port, opts...)
	if err != nil {
		log.Fatalln("fail to dial: ", err)
	}

	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)

	response, err := client.Rename(context.Background(), &pb.RenameRequest{oldPath, newPath})
	if err != nil {
		log.Fatalln("Rename Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
	}
	log.Println(response)
}
