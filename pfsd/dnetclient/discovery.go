package dnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/pnetclient"
	"github.com/cpssd/paranoid/pfsd/upnp"
	"log"
	"net"
	"time"
)

func SetDiscovery(host, port, serverPort string) {
	ipClient, _ := upnp.GetIP()
	ThisNode = globals.Node{
		IP:         ipClient,
		Port:       serverPort,
		CommonName: globals.CommonName,
	}
	globals.DiscoveryAddr = host + ":" + port

	if globals.TLSEnabled && !globals.TLSSkipVerify {
		// If host is an IP, we need to get the hostname via DNS
		ip := net.ParseIP(host)
		var dnsAddr string
		var err error
		if ip != nil {
			var dnsAddrs []string
			dnsAddrs, err = net.LookupAddr(ip.String())
			dnsAddr = dnsAddrs[0]
		} else {
			dnsAddr, err = net.LookupCNAME(host)
		}
		if err != nil {
			log.Println("ERROR: Could not complete DNS lookup:", err)
		}
		if dnsAddr == "" { // If no DNS entries exist
			log.Fatalln("FATAL: Can not find DNS entry for discovery server:", host)
		}
		discoveryCommonName = dnsAddr
	}
}

func JoinDiscovery(pool string) {
	if err := Join(pool); err != nil {
		if err = retryJoin(pool); err != nil {
			log.Fatalln("Failure dialing discovery server after multiple attempts, Giving up")
		}
	}
	globals.Wait.Add(2)
	go pingPeers()
	go renew()
}

// Periodically pings all known nodes on the network. Lives here and not
// in pnetclient since pnetclient is stateless and this function is more
// relevant to discovery.
func pingPeers() {
	// Ping as soon as this node joins
	pnetclient.Ping(globals.Nodes.GetAll())
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

func retryJoin(pool string) error {
	var err error
	for i := 0; i < 10; i++ {
		err = Join(pool)
		if err == nil {
			break
		}
	}
	return err
}

func renew() {
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
