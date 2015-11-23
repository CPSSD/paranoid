package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"log"
	"strconv"
)

func Utimes(ips []globals.Node, path,
	accessSeconds, accessMicroseconds, modifySeconds, modifyMicroseconds string) {

	for _, ipAddress := range ips {
		accessSecondsInt, err := strconv.ParseUint(accessSeconds, 10, 64)
		accessMicrosecondsInt, err := strconv.ParseUint(accessMicroseconds, 10, 64)
		modifySecondsInt, err := strconv.ParseUint(modifySeconds, 10, 64)
		modifyMicrosecondsInt, err := strconv.ParseUint(modifyMicroseconds, 10, 64)
		if err != nil {
			log.Println("Error parsing intergers.")
		}

		conn := Dial(ipAddress)

		defer conn.Close()
		client := pb.NewParanoidNetworkClient(conn)

		_, clientErr := client.Utimes(context.Background(),
			&pb.UtimesRequest{path, accessSecondsInt, accessMicrosecondsInt, modifySecondsInt, modifyMicrosecondsInt})
		if clientErr != nil {
			log.Println("Utimes Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
		}
	}
}
