package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	"google.golang.org/grpc"
	"log"
)

func Dial(ipAddress globals.Node) *grpc.ClientConn {
	conn, err := grpc.Dial(ipAddress.IP+":"+ipAddress.Port, grpc.WithInsecure())
	if err != nil {
		log.Println("fail to dial: ", err)
	}
	return conn
}
