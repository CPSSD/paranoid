package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
)

func Utimes(path string, accessSeconds, accessNanoSeconds, modifySeconds, modifyNanoSeconds int64) {
	nodes := globals.Nodes.GetAll()
	for _, node := range nodes {
		conn, err := Dial(node)
		if err != nil {
			Log.Error("Utimes: failed to dial ", node)
			continue
		}
		defer conn.Close()

		client := pb.NewParanoidNetworkClient(conn)

		_, clientErr := client.Utimes(context.Background(),
			&pb.UtimesRequest{path, accessSeconds, accessNanoSeconds, modifySeconds, modifyNanoSeconds})
		if clientErr != nil {
			Log.Error("Failed sending utimes to", node, "Error:", clientErr)
		}
	}
}
