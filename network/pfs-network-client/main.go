package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
)

var (
	dialer = websocket.Dialer{}
	pfsDir string
)

func main() {
	args := os.Args[1:]

	if len(args) != 4 {
		log.Fatalln("Usage: pfs-network-client --client <pfs-directory> <server-ip> <server-port>")
	}

	pfsDir = args[1]

	log.Println("Using pfs directory", pfsDir)

	serverAddr := args[2] + ":" + args[3]

	u := url.URL{Scheme: "ws", Host: serverAddr, Path: "/connect"}
	log.Println("Connecting to", u.String())

	c, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln("Cannot Connect:", err)
	}
	defer c.Close()

	go func() {
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Fatalln("Cannot Read Message:", err)
				break
			}
			parseMessage(message)
		}
	}()

}

func parseMessage(m []byte) {
	log.Println(m)
}
