package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/dnetclient"
	"github.com/cpssd/paranoid/pfsd/globals"
	"log"
	"time"
)

func SetDiscovery(ip, port, serverPort string) {
	log.Println(serverPort)
	ipClient, _ := GetIP()
	dnetclient.ThisNode = globals.Node{IP: ipClient, Port: serverPort}
	globals.DiscoveryAddr = ip + ":" + port
}

func JoinDiscovery(pool string) {
	if err := dnetclient.Join(pool); err != nil {
		connectionBuffer := 10
		log.Println("Error Connecting to Server, Attempting to reconnect")
		for connectionBuffer > 1 {
			err = dnetclient.Join(pool)
			connectionBuffer--
		}
		log.Println("Failure to connect to Discovery Server, Giving Up")
	} else {
		go renew()
	}
}

func renew() {
	for !globals.Disconnecting { //TODO change this to an async listener
		if err := dnetclient.Renew(); err != nil {
			log.Println("Failed to Renew Session")
		}
		globals.ResetInterval = 30000 // this is hard coded while I wait for interval fix
		time.Sleep(globals.ResetInterval * time.Millisecond)
	}
}

func Disconnect() {
	globals.Disconnecting = false
	dnetclient.Disconnect()
}
