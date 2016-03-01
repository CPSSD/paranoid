package pnetclient

import (
	"errors"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/upnp"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"strconv"
)

//Ping a peer asking to join raft network
func Ping() error {
	ip, err := upnp.GetIP()
	if err != nil {
		Log.Fatal("Can not ping peers: unable to get IP. Error:", err)
	}

	nodes := globals.Nodes.GetAll()
	for _, node := range nodes {
		Log.Info("Pinging node:", node)
		port := strconv.Itoa(globals.Port)

		conn, err := Dial(node)
		if err != nil {
			Log.Error("Ping: failed to dial ", node)
		}
		defer conn.Close()

		client := pb.NewParanoidNetworkClient(conn)

		_, err = client.Ping(context.Background(), &pb.PingRequest{
			Ip:         ip,
			Port:       port,
			CommonName: globals.CommonName,
			Uuid:       globals.UUID,
		})
		if err != nil {
			Log.Error("Error pinging", node, ":", err)
		} else {
			return nil
		}
	}
	return errors.New("unable to join raft network, no peer has returned okay")
}
