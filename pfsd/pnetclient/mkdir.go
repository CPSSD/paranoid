package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func Mkdir(directory string, mode uint32) {
	nodes := globals.Nodes.GetAll()
	for _, node := range nodes {
		conn, err := Dial(node)
		if err != nil {
			log.Println("Mkdir error failed to dial ", node)
			continue
		}
		defer conn.Close()

		client := pb.NewParanoidNetworkClient(conn)

		_, err = client.Mkdir(context.Background(), &pb.MkdirRequest{directory, mode})
		if err != nil {
			log.Println("Mkdir Error on ", node, "Error:", err)
		}
	}
}
