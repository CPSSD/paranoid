package network

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

func ping(ips []globals.Node) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	for _, value := range ips {
		pingServer(value, opts)
	}
}

func pingServer(ipAddress globals.Node, opts []grpc.DialOption) {
	ip := GetIP()
	conn, err := grpc.Dial(ipAddress.IP+":"+ipAddress.Port, opts...)
	if err != nil {
		log.Fatalln("fail to dial: ", err)
	}

	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)

	port := strconv.Itoa(globals.Port)
	response, err := client.Ping(context.Background(), &pb.PingRequest{ip, port})
	if err != nil {
		log.Fatalln("Cant Ping ", ipAddress.IP+":"+ipAddress.Port)
	}
	log.Println(response)
}
