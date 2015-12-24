package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func Link(oldPath, newPath string) {
	nodes := globals.Nodes.GetAll()
	for _, node := range nodes {
		conn, err := Dial(node)
		if err != nil {
			log.Println("Link error failed to dial ", node)
			continue
		}
		defer conn.Close()

		client := pb.NewParanoidNetworkClient(conn)

		_, err = client.Link(context.Background(), &pb.LinkRequest{oldPath, newPath})
		if err != nil {
			log.Println("Link Error on ", node, "Error:", err)
		}
	}
}
