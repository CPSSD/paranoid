package network

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

func write(ips []globals.Node, path string, data []byte, offset, length string) {

	for _, ipAddress := range ips {
		response, err := sendWriteRequest(ipAddress, path, data, offset, length)
		if err != nil {
			log.Println("Error Sending Message: ", err)
		} else {
			log.Println(response)
		}
	}
}

func sendWriteRequest(ipAddress globals.Node, path string, data []byte, offset, length string) (string, error) {

	conn, err := grpc.Dial(ipAddress.IP+":"+ipAddress.Port, grpc.WithInsecure())
	if err != nil {
		log.Println("fail to dial: ", err)
		return "", err
	}

	defer conn.Close()
	client := pb.NewParanoidNetworkClient(conn)

	offsetInt, _ := strconv.ParseUint(offset, 10, 64)
	lengthInt, _ := strconv.ParseUint(length, 10, 64)
	response, err := client.Write(context.Background(), &pb.WriteRequest{path, data, offsetInt, lengthInt})
	if err != nil {
		log.Println("Write Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)

	}
	return response.String(), err
}
