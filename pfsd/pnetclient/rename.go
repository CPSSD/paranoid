package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func rename(ips []globals.Node, oldPath, newPath string) {
	for _, ipAddress := range ips {
		conn := Dial(ipAddress)

		defer conn.Close()
		client := pb.NewParanoidNetworkClient(conn)

		_, err := client.Rename(context.Background(), &pb.RenameRequest{oldPath, newPath})
		if err != nil {
			log.Println("Rename Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
		}
	}
}
