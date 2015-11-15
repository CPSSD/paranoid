package network

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

func truncate(ips []globals.Node, path string, length string) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	for _, ipAddress := range ips {
		sendTruncateRequest(ipAddress, path, length, opts)
	}
}

func sendTruncateRequest(ipAddress globals.Node, path string, length string, opts []grpc.DialOption) {
	conn, err := grpc.Dial(ipAddress.IP+":"+ipAddress.Port, opts...)
	if err != nil {
		log.Fatalln("fail to dial: ", err)
	}

	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)
	lengthInt, _ := strconv.ParseUint(length, 10, 64)

	response, err := client.Truncate(context.Background(), &pb.TruncateRequest{path, lengthInt})
	if err != nil {
		log.Fatalln("Truncate Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
	}
	log.Println(response)
}
