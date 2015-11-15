package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

func utimes(ips []globals.Node, path,
	accessSeconds, accessMicroseconds, modifySeconds, modifyMicroseconds string) {

	for _, ipAddress := range ips {
		sendUtimesRequest(ipAddress, path, accessSeconds,
			accessMicroseconds, modifySeconds, modifyMicroseconds)
	}
}

func sendUtimesRequest(ipAddress globals.Node, path,
	accessSeconds, accessMicroseconds, modifySeconds, modifyMicroseconds string) {

	accessSecondsInt, _ := strconv.ParseUint(accessSeconds, 10, 64)
	accessMicrosecondsInt, _ := strconv.ParseUint(accessMicroseconds, 10, 64)
	modifySecondsInt, _ := strconv.ParseUint(modifySeconds, 10, 64)
	modifyMicrosecondsInt, _ := strconv.ParseUint(modifyMicroseconds, 10, 64)

	conn, err := grpc.Dial(ipAddress.IP+":"+ipAddress.Port, grpc.WithInsecure())
	if err != nil {
		log.Println("fail to dial: ", err)
	}

	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)

	response, err := client.Utimes(context.Background(),
		&pb.UtimesRequest{path, accessSecondsInt, accessMicrosecondsInt, modifySecondsInt, modifyMicrosecondsInt})
	if err != nil {
		log.Println("Utimes Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
	}
	log.Println(response)
}
