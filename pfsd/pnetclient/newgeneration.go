package pnetclient

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
)

// NewGeneration is used to create a new KeyPair generation in the cluster,
// prior to this node joining.
func NewGeneration(password string) error {
	nodes := globals.Nodes.GetAll()
	for _, node := range nodes {
		Log.Info("Sending new generation request to node:", node)

		conn, err := Dial(node)
		if err != nil {
			Log.Error("NewGeneration: failed to dial", node)
		}
		defer conn.Close()

		client := pb.NewParanoidNetworkClient(conn)

		resp, err := client.NewGeneration(context.Background(), &pb.NewGenerationRequest{
			RequestingNode: &pb.Node{
				Ip:         globals.ThisNode.IP,
				Port:       globals.ThisNode.Port,
				CommonName: globals.ThisNode.CommonName,
				Uuid:       globals.ThisNode.UUID,
			},
			PoolPassword: password,
		})
		if err != nil {
			Log.Error("Error requesting to create new generation", node, ":", err)
		} else {
			peers := make([]globals.Node, len(resp.Peers))
			for i, v := range resp.Peers {
				peers[i] = globals.Node{
					IP:         v.Ip,
					Port:       v.Port,
					CommonName: v.CommonName,
					UUID:       v.Uuid,
				}
			}
			err := Distribute(globals.EncryptionKey, peers, int(resp.GenerationNumber))
			if err != nil {
				Log.Error("Failed to distribute key:", err)
				return fmt.Errorf("failed to distribute key: %s", err)
			}
		}
	}
	return errors.New("unable to create new generation, no peer has returned okay")
}
