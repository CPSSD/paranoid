package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

func truncate(ips []globals.Node, path string, length string) {
	for _, ipAddress := range ips {
		sendTruncateRequest(ipAddress, path, length)
	}
}

func sendTruncateRequest(ipAddress globals.Node, path string, length string) {
	conn, err := grpc.Dial(ipAddress.IP+":"+ipAddress.Port, grpc.WithInsecure())
	if err != nil {
		log.Println("fail to dial: ", err)
	}

	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	lengthInt, _ := strconv.ParseUint(length, 10, 64)

	response, err := client.Truncate(context.Background(), &pb.TruncateRequest{path, lengthInt})
	if err != nil {
		log.Println("Truncate Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
	}
	log.Println(response)
}
