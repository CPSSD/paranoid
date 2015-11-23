package dnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/pnetclient"
	"log"
	"time"
)

func SetDiscovery(ip, port, serverPort string) {
	ipClient, _ := pnetclient.GetIP()
	ThisNode = globals.Node{IP: ipClient, Port: serverPort}
	globals.DiscoveryAddr = ip + ":" + port
}

func JoinDiscovery(pool string) {
	if err := Join(pool); err != nil {
		connectionBuffer := 10
		log.Println("Error Connecting to Server, Attempting to reconnect")
		for connectionBuffer > 1 {
			err = Join(pool)
			connectionBuffer--
		}
		log.Println("Failure to connect to Discovery Server, Giving Up")
	} else {
		go renew()
	}
}

func renew() {
	for {
		if err := Renew(); err != nil {
			log.Println("Failed to Renew Session")
		}
		globals.ResetInterval = 30000 // this is hard coded while I wait for interval fix
		time.Sleep(globals.ResetInterval * time.Millisecond)
	}
}
