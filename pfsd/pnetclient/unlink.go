package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func unlink(ips []globals.Node, path string) {
	for _, ipAddress := range ips {
		conn := Dial(ipAddress)

		defer conn.Close()
		client := pb.NewParanoidNetworkClient(conn)

		_, err := client.Unlink(context.Background(), &pb.UnlinkRequest{path})
		if err != nil {
			log.Println("Unlink Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
		}
	}
}
