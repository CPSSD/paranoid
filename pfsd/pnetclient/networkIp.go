package pnetclient

import (
	"errors"
	"net"
)

func GetIP() (string, error) {
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("No IP found")
}
