package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func Truncate(path string, length uint64) {
	nodes := globals.Nodes.GetAll()
	for _, node := range nodes {
		conn, err := Dial(node)
		if err != nil {
			log.Println("Truncate error failed to dial ", node)
			continue
		}
		defer conn.Close()

		client := pb.NewParanoidNetworkClient(conn)

		_, clientErr := client.Truncate(context.Background(), &pb.TruncateRequest{path, length})
		if clientErr != nil {
			log.Println("Truncate Error on ", node, "Error:", err)
		}
	}
}
