package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
	"strconv"
)

func Creat(ips []globals.Node, filename, permissions string) {
	for _, ipAddress := range ips {
		permissions64, _ := strconv.ParseUint(permissions, 8, 32)
		permissionsInt := uint32(permissions64)

		conn := Dial(ipAddress)

		defer conn.Close()
		client := pb.NewParanoidNetworkClient(conn)

		_, err := client.Creat(context.Background(), &pb.CreatRequest{filename, permissionsInt})
		if err != nil {
			log.Println("Failure Sending Message to", ipAddress, " Error:", err)
		}
	}
}
