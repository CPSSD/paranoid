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
		accessSecondsInt, _ := strconv.ParseUint(accessSeconds, 10, 64)
		accessMicrosecondsInt, _ := strconv.ParseUint(accessMicroseconds, 10, 64)
		modifySecondsInt, _ := strconv.ParseUint(modifySeconds, 10, 64)
		modifyMicrosecondsInt, _ := strconv.ParseUint(modifyMicroseconds, 10, 64)

		conn := Dial(ipAddress)

		defer conn.Close()
		client := pb.NewParanoidNetworkClient(conn)

		_, err := client.Utimes(context.Background(),
			&pb.UtimesRequest{path, accessSecondsInt, accessMicrosecondsInt, modifySecondsInt, modifyMicrosecondsInt})
		if err != nil {
			log.Println("Utimes Error on ", ipAddress.IP+":"+ipAddress.Port, "Error:", err)
		}
	}
}
