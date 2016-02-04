package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
)

func Link(oldPath, newPath string) {
	nodes := globals.Nodes.GetAll()
	for _, node := range nodes {
		conn, err := Dial(node)
		if err != nil {
			Log.Error("Link: Failed to dial ", node)
			continue
		}
		defer conn.Close()

		client := pb.NewParanoidNetworkClient(conn)

		_, err = client.Link(context.Background(), &pb.LinkRequest{oldPath, newPath})
		if err != nil {
			Log.Error("Failure sending link to", node, "Error:", err)
		}
	}
}
