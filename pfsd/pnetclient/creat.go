package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
)

func Creat(filename string, permissions uint32) {
	nodes := globals.Nodes.GetAll()
	for _, node := range nodes {
		conn, err := Dial(node)
		if err != nil {
			Log.Error("Creat: failed to dial ", node)
			continue
		}
		defer conn.Close()

		client := pb.NewParanoidNetworkClient(conn)

		_, err = client.Creat(context.Background(), &pb.CreatRequest{filename, permissions})
		if err != nil {
			Log.Error("Failure sending creat to", node, " Error:", err)
		}
	}
}
