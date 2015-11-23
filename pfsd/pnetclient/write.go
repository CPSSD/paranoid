package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
	"strconv"
)

func Write(ips []globals.Node, path string, data []byte, offset, length string) {
	for _, ipAddress := range ips {
		conn := Dial(ipAddress)

		defer conn.Close()
		client := pb.NewParanoidNetworkClient(conn)

		offsetInt, _ := strconv.ParseUint(offset, 10, 64)
		lengthInt, _ := strconv.ParseUint(length, 10, 64)
		_, err := client.Write(context.Background(), &pb.WriteRequest{path, data, offsetInt, lengthInt})
		if err != nil {
			log.Println("Write Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
		}
	}
}
