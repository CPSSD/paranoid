package pnetclient

import (
	"errors"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/upnp"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
)

//JoinCluster is used to request to join a raft cluster
func JoinCluster(password string) error {
	ip, err := upnp.GetIP()
	if err != nil {
		Log.Fatal("Can not contact peers: unable to get IP. Error:", err)
	}

	nodes := globals.Nodes.GetAll()
	for _, node := range nodes {
		Log.Info("Sending join cluster request to node:", node)

		conn, err := Dial(node)
		if err != nil {
			Log.Error("JoinCluster: failed to dial ", node)
		}
		defer conn.Close()

		client := pb.NewParanoidNetworkClient(conn)

		_, err = client.JoinCluster(context.Background(), &pb.JoinClusterRequest{
			Ip:           ip,
			Port:         globals.ThisNode.Port,
			CommonName:   globals.ThisNode.CommonName,
			Uuid:         globals.ThisNode.UUID,
			PoolPassword: password,
		})
		if err != nil {
			Log.Error("Error requesting to join cluster", node, ":", err)
		} else {
			return nil
		}
	}
	return errors.New("unable to join raft network, no peer has returned okay")
}
