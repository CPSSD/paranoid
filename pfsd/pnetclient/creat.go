package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func Creat(filename string, permissions uint32) {
	nodes := globals.Nodes.GetAll()
	for _, node := range nodes {
		conn, err := Dial(node)
		if err != nil {
			log.Println("Creat error failed to dial ", node)
			continue
		}
		defer conn.Close()

		client := pb.NewParanoidNetworkClient(conn)

		_, err = client.Creat(context.Background(), &pb.CreatRequest{filename, permissions})
		if err != nil {
			log.Println("Failure Sending Message to", node, " Error:", err)
		}
	}
}
