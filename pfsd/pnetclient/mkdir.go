package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
	"strconv"
)

// Mkdir is used to create directories
func Mkdir(ips []globals.Node, directory string, mode string) {
	for _, ipAddress := range ips {
		conn := Dial(ipAddress)

		intMode, err := strconv.ParseUint(mode, 8, 32)
		if err != nil {
			log.Println("Error parsing mode in Mkdir.")
		}

		defer conn.Close()
		client := pb.NewParanoidNetworkClient(conn)

		_, err = client.Mkdir(context.Background(), &pb.MkdirRequest{directory, uint32(intMode)})
		if err != nil {
			log.Println("Mkdir Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
		}
	}
}
