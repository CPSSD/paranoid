package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
	"os"
	"strconv"
)

func Ping(ips []globals.Node) {
	for _, ipAddress := range ips {
		ip, _ := GetIP()
		port := strconv.Itoa(globals.Port)
		hostname, err := os.Hostname()
		if err != nil {
			log.Println("ERROR: Could not get machine hostname:", err)
		}

		conn := Dial(ipAddress)
		defer conn.Close()
		client := pb.NewParanoidNetworkClient(conn)

		_, err = client.Ping(context.Background(), &pb.PingRequest{ip, port, hostname})
		if err != nil {
			log.Println("Can't Ping ", ipAddress.IP+":"+ipAddress.Port)
		}
	}
}
