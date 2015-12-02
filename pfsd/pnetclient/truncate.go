package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
	"strconv"
)

func Truncate(ips []globals.Node, path string, length string) {
	for _, ipAddress := range ips {
		conn := Dial(ipAddress)

		defer conn.Close()
		client := pb.NewParanoidNetworkClient(conn)
		lengthInt, err := strconv.ParseUint(length, 10, 64)
		if err != nil {
			log.Println("Error parsing intergers.")
		}

		_, clientErr := client.Truncate(context.Background(), &pb.TruncateRequest{path, lengthInt})
		if clientErr != nil {
			log.Println("Truncate Error on ", ipAddress, "Error:", err)
		}
	}
}
