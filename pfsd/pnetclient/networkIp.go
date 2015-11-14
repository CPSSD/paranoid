package network

import (
	"net"
)

func GetIP() string {
	ifaces, _ := net.Interfaces()
	// handle err
	addrs, _ := ifaces[1].Addrs()
	var ip net.IP
	switch v := addrs[0].(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	return ip.String()
}
