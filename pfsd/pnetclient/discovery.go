package network

import (
	"github.com/cpssd/paranoid/pfsd/dnetclient"
	"github.com/cpssd/paranoid/pfsd/globals"
	"log"
	"time"
)

func SetDiscovery(ip, port string) {
	dnetclient.ThisNode = globals.Node{IP: ip, Port: port}
	globals.DiscoveryAddr = ip + ":" + port
}

func Join(pool string) {
	dnetclient.Join(pool)
	if err := dnetclient.Join(pool); err != nil {
		connectionBuffer := 30
		// Going to attempt connecting 30 times before I give up
		for connectionBuffer > 1 {
			err = dnetclient.Join(pool)
		}
	} else {
		go renew()
	}
}

func renew() {
	for true { //noway to break this except forcefully
		time.Sleep(globals.ResetInterval)
		err := dnetclient.Renew()
		if err != nil {
			log.Println("failure to renew connection")
		}
	}
}
