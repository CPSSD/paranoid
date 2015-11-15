package network

import (
	"errors"
	"net"
	"strings"
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
			if strings.Contains(ip.String(), "192.168.") {
				return ip.String(), nil
			}
			if strings.Contains(ip.String(), "10.") {
				return ip.String(), nil
			}
			if strings.Contains(ip.String(), "17.16.") {
				return ip.String(), nil
			}
		}
	}
	return "", errors.New("No IP found")
}
