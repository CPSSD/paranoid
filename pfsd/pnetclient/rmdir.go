package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

// Rmdir is used to create directories
func Rmdir(ips []globals.Node, directory string) {
	for _, ipAddress := range ips {
		conn := Dial(ipAddress)

		defer conn.Close()
		client := pb.NewParanoidNetworkClient(conn)

		_, err := client.Rmdir(context.Background(), &pb.RmdirRequest{directory})
		if err != nil {
			log.Println("Rmdir Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
		}
	}
}
