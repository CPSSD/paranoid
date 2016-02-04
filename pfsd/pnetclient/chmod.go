package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
)

func Chmod(path string, mode uint32) {
	nodes := globals.Nodes.GetAll()
	for _, node := range nodes {
		conn, err := Dial(node)
		if err != nil {
			Log.Error("Chmod: Failed to dial ", node)
			continue
		}
		defer conn.Close()

		client := pb.NewParanoidNetworkClient(conn)
		_, err = client.Chmod(context.Background(), &pb.ChmodRequest{path, mode})
		if err != nil {
			Log.Error("Failure sending chmod to", node, "Error:", err)
		}
	}
}
