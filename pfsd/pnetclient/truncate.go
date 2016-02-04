package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
)

func Truncate(path string, length uint64) {
	nodes := globals.Nodes.GetAll()
	for _, node := range nodes {
		conn, err := Dial(node)
		if err != nil {
			Log.Error("Truncate: failed to dial ", node)
			continue
		}
		defer conn.Close()

		client := pb.NewParanoidNetworkClient(conn)

		_, clientErr := client.Truncate(context.Background(), &pb.TruncateRequest{path, length})
		if clientErr != nil {
			Log.Error("Failed sending truncate to", node, "Error:", err)
		}
	}
}
