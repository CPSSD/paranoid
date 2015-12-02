package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

// Mkdir is used to create directories
func Mkdir(ips []globals.Node, directory string, mode uint32) {
	for _, ipAddress := range ips {
		conn := Dial(ipAddress)

		defer conn.Close()
		client := pb.NewParanoidNetworkClient(conn)

		_, err := client.Mkdir(context.Background(), &pb.MkdirRequest{directory, mode})
		if err != nil {
			log.Println("Mkdir Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
		}
	}
}
