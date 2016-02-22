package dnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/pnetclient"
	"github.com/cpssd/paranoid/pfsd/upnp"
	"net"
	"time"
)

func SetDiscovery(host, port, serverPort string) {
	ipClient, _ := upnp.GetIP()
	globals.ThisNode = globals.Node{
		IP:         ipClient,
		Port:       serverPort,
		CommonName: globals.CommonName,
		UUID:       globals.UUID,
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
			Log.Error("Could not complete DNS lookup:", err)
		}
		if dnsAddr == "" { // If no DNS entries exist
			Log.Fatal("Can not find DNS entry for discovery server:", host)
		}
		discoveryCommonName = dnsAddr
	}
}

func JoinDiscovery(pool string) {
	if err := Join(pool); err != nil {
		if err = retryJoin(pool); err != nil {
			Log.Fatal("Failure dialing discovery server after multiple attempts, Giving up", err)
		}
	}
	globals.Wait.Add(1)
	go pingPeers()
}

// Periodically pings all known nodes on the network. Lives here and not
// in pnetclient since pnetclient is stateless and this function is more
// relevant to discovery.
func pingPeers() {
	// Ping as soon as this node joins
	pnetclient.Ping()
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
			pnetclient.Ping()
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
