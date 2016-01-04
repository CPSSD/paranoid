package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
)

func Utimes(path string, accessSeconds, accessNanoSeconds, modifySeconds, modifyNanoSeconds int64) {
	nodes := globals.Nodes.GetAll()
	for _, node := range nodes {
		conn, err := Dial(node)
		if err != nil {
			log.Println("Utimes error failed to dial ", node)
			continue
		}
		defer conn.Close()

		client := pb.NewParanoidNetworkClient(conn)

		_, clientErr := client.Utimes(context.Background(),
			&pb.UtimesRequest{path, accessSeconds, accessNanoSeconds, modifySeconds, modifyNanoSeconds})
		if clientErr != nil {
			log.Println("Utimes Error on ", node, "Error:", clientErr)
		}
	}
}
