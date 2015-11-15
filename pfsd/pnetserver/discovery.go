package pnetserver

import (
	"github.com/cpssd/paranoid/pfsd/dnetclient"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/pnetclient"
	"log"
	"time"
)

func SetDiscovery(ip, port, serverPort string) {
	log.Println(serverPort)
	ipClient, _ := pnetclient.GetIP()
	dnetclient.ThisNode = globals.Node{IP: ipClient, Port: serverPort}
	globals.DiscoveryAddr = ip + ":" + port
}

func JoinDiscovery(pool string) {
	if err := dnetclient.Join(pool); err != nil {
		log.Println("I'm an error ")
		connectionBuffer := 10
		log.Println("Error Connecting to Server, Attempting to reconnect")
		for connectionBuffer > 1 {
			err = dnetclient.Join(pool)
			connectionBuffer--
		}
	} else {
		go renew()
	}
}

func renew() {
	for { //Cant be terminated right now. Going to write a call to check if Disconnect has been called
		if err := dnetclient.Renew(); err != nil {
			log.Println("Failed to Renew Session")
		}
		globals.ResetInterval = 5000 // this is hard coded while I wait for interval fix
		time.Sleep(globals.ResetInterval * time.Millisecond)
	}
}
