package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
	"strconv"
)

func chmod(ips []globals.Node, path string, mode string) {
	for _, ipAddress := range ips {
		mode64, _ := strconv.ParseUint(mode, 8, 32)
		modeInt := uint32(mode64)

		conn := Dial(ipAddress)

		defer conn.Close()
		client := pb.NewParanoidNetworkClient(conn)

		_, err := client.Chmod(context.Background(), &pb.ChmodRequest{path, modeInt})
		if err != nil {
			log.Println("Chmod Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
		}
	}
}
