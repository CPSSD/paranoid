package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
	"strconv"
)

func truncate(ips []globals.Node, path string, length string) {
	for _, ipAddress := range ips {
		conn := Dial(ipAddress)

		defer conn.Close()
		client := pb.NewParanoidNetworkClient(conn)
		lengthInt, _ := strconv.ParseUint(length, 10, 64)

		_, err := client.Truncate(context.Background(), &pb.TruncateRequest{path, lengthInt})
		if err != nil {
			log.Println("Truncate Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
		}
	}
}
