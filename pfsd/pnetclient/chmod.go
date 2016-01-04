package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func Chmod(path string, mode uint32) {
	nodes := globals.Nodes.GetAll()
	for _, node := range nodes {
		conn, err := Dial(node)
		if err != nil {
			log.Println("Chmod error failed to dial ", node)
			continue
		}
		defer conn.Close()

		client := pb.NewParanoidNetworkClient(conn)
		_, err = client.Chmod(context.Background(), &pb.ChmodRequest{path, mode})
		if err != nil {
			log.Println("Chmod Error on ", node, "Error:", err)
		}
	}
}
