package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/upnp"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
	"strconv"
)

func Ping(ips []globals.Node) {
	ip, _ := upnp.GetIP()
	for _, ipAddress := range ips {
		conn := Dial(ipAddress)

		defer conn.Close()
		client := pb.NewParanoidNetworkClient(conn)

		port := strconv.Itoa(globals.Port)
		_, err := client.Ping(context.Background(), &pb.PingRequest{ip, port})
		if err != nil {
			log.Println("Can't Ping ", ipAddress.IP+":"+ipAddress.Port)
		}
	}
}
