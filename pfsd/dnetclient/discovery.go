package dnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/pnetclient"
	"log"
	"time"
)

const peerPingInterval time.Duration = time.Minute

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
		return
	}
	globals.Wait.Add(2)
	go pingPeers()
	go renew()
}

// Periodically pings all known nodes on the network. Lives here and not
// in pnetclient since pnetclient is stateless and this function is more
// relevant to discovery.
func pingPeers() {
	timer := time.NewTimer(peerPingInterval)
	defer timer.Stop()
	defer globals.Wait.Done()
	for {
		select {
		case _, ok := <-globals.Quit:
			if !ok {
				return
			}
		case <-timer.C:
			pnetclient.Ping(globals.Nodes.GetAll())
			timer.Reset(peerPingInterval)
		}
	}
}

func renew() {
	// TODO(sean) Set this to the actual reset interval when implemented
	globals.ResetInterval = 30000
	timer := time.NewTimer(globals.ResetInterval * time.Millisecond)
	defer timer.Stop()
	defer globals.Wait.Done()
	for {
		select {
		case _, ok := <-globals.Quit:
			if !ok {
				log.Println("INFO: Disconnected from discovery server.")
				return
			}
		case <-timer.C:
			if err := Renew(); err != nil {
				log.Println("Failed to Renew Session")
			}
			timer.Reset(globals.ResetInterval * time.Millisecond)
		}
	}
}
