package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

func chmod(ips []globals.Node, path string, mode string) {
	for _, ipAddress := range ips {
		sendChmodRequest(ipAddress, path, mode)
	}
}

func sendChmodRequest(ipAddress globals.Node, path, mode string) {
	var modeInt uint32
	mode64, _ := strconv.ParseUint(mode, 8, 32)
	modeInt = uint32(mode64)

	conn, err := grpc.Dial(ipAddress.IP+":"+ipAddress.Port, grpc.WithInsecure())
	if err != nil {
		log.Fatalln("fail to dial: ", err)
	}

	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)

	response, err := client.Chmod(context.Background(), &pb.ChmodRequest{path, modeInt})
	if err != nil {
		log.Println("Chmod Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
	}
	log.Println(response)
}
