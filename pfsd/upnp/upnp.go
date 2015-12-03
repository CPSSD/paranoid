package upnp

import (
	"errors"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/huin/goupnp/dcps/internetgateway1"
	"math/rand"
	"net"
)

var (
	uPnPipClients       []*internetgateway1.WANIPConnection1
	uPnPpppClients      []*internetgateway1.WANPPPConnection1
	ipPortMappedClient  *internetgateway1.WANIPConnection1
	pppPortMappedClient *internetgateway1.WANPPPConnection1
)

const attemptedPortAssignments = 10

//Discovers uPnP devices on the network.
func DiscoverDevices() error {
	ipclients, _, err := internetgateway1.NewWANIPConnection1Clients()
	if err == nil {
		uPnPipClients = ipclients
	}
	pppclients, _, err := internetgateway1.NewWANPPPConnection1Clients()
	if err == nil {
		uPnPpppClients = pppclients
	}
	if len(ipclients) > 0 || len(pppclients) > 0 {
		return nil
	}
	return errors.New("No devices found")
}

func getUnOccupiedPortsIp(client *internetgateway1.WANIPConnection1) []int {
	m := make(map[int]bool)
	i := 0
	for {
		_, port, _, _, _, _, _, _, err := client.GetGenericPortMappingEntry(uint16(i))
		if err != nil {
			break
		}
		i++
		m[int(port)] = true
	}
	openPorts := make([]int, 0)
	for i = 1; i < 65536; i++ {
		if m[i] == false {
			openPorts = append(openPorts, i)
		}
	}
	return openPorts
}

func getUnOccupiedPortsppp(client *internetgateway1.WANPPPConnection1) []int {
	m := make(map[int]bool)
	i := 0
	for {
		_, port, _, _, _, _, _, _, err := client.GetGenericPortMappingEntry(uint16(i))
		if err != nil {
			break
		}
		i++
		m[int(port)] = true
	}
	openPorts := make([]int, 0)
	for i = 1; i < 65536; i++ {
		if m[i] == false {
			openPorts = append(openPorts, i)
		}
	}
	return openPorts
}

func AddPortMapping(internalPort int) (int, error) {
	ip, err := GetInternalIp()
	if err != nil {
		return 0, err
	}
	for _, client := range uPnPipClients {
		openPorts := getUnOccupiedPortsIp(client)
		if len(openPorts) > 0 {
			for i := 0; i < attemptedPortAssignments; i++ {
				port := openPorts[rand.Intn(len(openPorts)-1)]
				err := client.AddPortMapping("", uint16(internalPort), "TCP", uint16(port), ip, true, "", 0)
				if err == nil {
					ipPortMappedClient = client
					return port, nil
				}
			}
		}
	}
	for _, client := range uPnPpppClients {
		openPorts := getUnOccupiedPortsppp(client)
		if len(openPorts) > 0 {
			for i := 0; i < attemptedPortAssignments; i++ {
				port := openPorts[rand.Intn(len(openPorts)-1)]
				err := client.AddPortMapping("", uint16(internalPort), "TCP", uint16(port), ip, true, "", 0)
				if err == nil {
					pppPortMappedClient = client
					return port, nil
				}
			}
		}
	}
	return 0, errors.New("Unable to map port")
}

func ClearPortMapping(externalPort int) error {
	if ipPortMappedClient != nil {
		return ipPortMappedClient.DeletePortMapping("", uint16(externalPort), "TCP")
	}
	if pppPortMappedClient == nil {
		return pppPortMappedClient.DeletePortMapping("", uint16(externalPort), "TCP")
	}
	return errors.New("No uPnP device available")
}

//Gets the internal Ip address
func GetInternalIp() (string, error) {
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

//Gets the external IP of the port mapped device.
func GetExternalIp() (string, error) {
	if ipPortMappedClient != nil {
		externalIp, err := ipPortMappedClient.GetExternalIPAddress()
		if err == nil {
			return externalIp, nil
		}
	}
	if pppPortMappedClient != nil {
		externalIp, err := pppPortMappedClient.GetExternalIPAddress()
		if err == nil {
			return externalIp, nil
		}
	}
	return "", errors.New("Unabled to get get external IP address")
}

func GetIP() (string, error) {
	if globals.UpnpEnabled {
		return GetExternalIp()
	} else {
		return GetInternalIp()
	}
}
