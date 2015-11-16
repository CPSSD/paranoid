package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

func ping(ips []globals.Node) {
	for _, value := range ips {
		pingServer(value)
	}
}

func pingServer(ipAddress globals.Node) {
	ip, _ := GetIP()
	conn, err := grpc.Dial(ipAddress.IP+":"+ipAddress.Port, grpc.WithInsecure())
	if err != nil {
		log.Println("fail to dial: ", err)
	}

	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)

	port := strconv.Itoa(globals.Port)
	response, err := client.Ping(context.Background(), &pb.PingRequest{ip, port})
	if err != nil {
		log.Println("Cant Ping ", ipAddress.IP+":"+ipAddress.Port)
	}
	log.Println(response)
}
