package network

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

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	for _, ipAddress := range ips {
		sendUtimesRequest(ipAddress, path, accessSeconds,
			accessMicroseconds, modifySeconds, modifyMicroseconds, opts)
	}
}

func sendUtimesRequest(ipAddress globals.Node, path,
	accessSeconds, accessMicroseconds, modifySeconds, modifyMicroseconds string, opts []grpc.DialOption) {

	accessSecondsInt, _ := strconv.ParseUint(accessSeconds, 10, 64)
	accessMicrosecondsInt, _ := strconv.ParseUint(accessMicroseconds, 10, 64)
	modifySecondsInt, _ := strconv.ParseUint(modifySeconds, 10, 64)
	modifyMicrosecondsInt, _ := strconv.ParseUint(modifyMicroseconds, 10, 64)

	conn, err := grpc.Dial(ipAddress.IP+":"+ipAddress.Port, opts...)
	if err != nil {
		log.Fatalln("fail to dial: ", err)
	}

	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)

	response, err := client.Utimes(context.Background(),
		&pb.UtimesRequest{path, accessSecondsInt, accessMicrosecondsInt, modifySecondsInt, modifyMicrosecondsInt})
	if err != nil {
		log.Fatalln("Utimes Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
	}
	log.Println(response)
}
