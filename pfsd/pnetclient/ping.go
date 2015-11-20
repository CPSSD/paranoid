package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
	"strconv"
)

func ping(ips []globals.Node) {
	for _, ipAddress := range ips {
		ip, _ := GetIP()
		conn := Dial(ipAddress)

		defer conn.Close()
		client := pb.NewParanoidNetworkClient(conn)

		port := strconv.Itoa(globals.Port)
		_, err := client.Ping(context.Background(), &pb.PingRequest{ip, port})
		if err != nil {
			log.Println("Cant Ping ", ipAddress.IP+":"+ipAddress.Port)
		}
	}
}

func pingServer(ipAddress globals.Node) {

}
