package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/upnp"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
	"strconv"
)

func Ping() {
	ip, err := upnp.GetIP()
	if err != nil {
		log.Fatalln("Can not ping peers: unabled to get IP. Error:", err)
	}

	nodes := globals.Nodes.GetAll()
	for _, node := range nodes {
		port := strconv.Itoa(globals.Port)

		conn, err := Dial(node)
		if err != nil {
			log.Println("Ping error: failed to dial ", node)
		}
		defer conn.Close()

		client := pb.NewParanoidNetworkClient(conn)

		_, err = client.Ping(context.Background(), &pb.PingRequest{ip, port, globals.CommonName})
		if err != nil {
			log.Println("Can't Ping ", node)
		}
	}
}
