package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
)

func GetIP() (string, error) {
	err := globals.UpnpMapping.SearchGateway()
	if err == nil {
		//Upnp supported
		err := globals.UpnpMapping.ExternalIPAddr()
		if err != nil {
			return "", err
		}
		return globals.UpnpMapping.GatewayOutsideIP, nil
	} else {
		return globals.UpnpMapping.LocalHost, nil
	}
}
