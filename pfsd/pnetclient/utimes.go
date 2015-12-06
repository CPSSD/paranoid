package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func Utimes(ips []globals.Node, path string, data []byte) {
	for _, ipAddress := range ips {
		conn := Dial(ipAddress)

		defer conn.Close()
		client := pb.NewParanoidNetworkClient(conn)

		_, clientErr := client.Utimes(context.Background(),
			&pb.UtimesRequest{path, data})
		if clientErr != nil {
			log.Println("Utimes Error on ", ipAddress, "Error:", clientErr)
		}
	}
}
