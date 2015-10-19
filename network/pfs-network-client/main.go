package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"strconv"
)

var (
	pfsDir string
)

type uuid string
type base64 []byte

type MessageData struct {
	sender     uuid
	actionType string
	name       string
	target     string
	offset     string
	length     string
	data       base64
}

func main() {
	args := os.Args[1:]

	if len(args) != 4 {
		fmt.Println("Usage:\n\tpfs-network-client --client <pfs-directory> <server-ip> <server-port>")
		os.Exit(1)
	}

	pfsDir = args[1]

	log.Println("Using pfs directory", pfsDir)

	address := args[2]
	port, err := strconv.Atoi(args[3]) // Get the port and parse it to int

	if err != nil || port < 1 || port > 65545 {
		log.Fatalln("FATAL: port must be a number between 1 and 65535, inclusive")
	}

	serverAddr := address + ":" + strconv.Itoa(port)

	u := url.URL{Scheme: "ws", Host: serverAddr, Path: "/connect"}
	log.Println("Connecting to server...")

	dialer := websocket.Dialer{}
	c, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalln("FATAL: Cannot Connect:", err)
	} else {
		log.Println("Establised connection with", u.String())
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatalln("Cannot Read Message:", err)
			break
		}
		parseMessage(message)
	}

	defer c.Close()
}

func parseMessage(messageString []byte) {
	m := MessageData{}

	if err := json.Unmarshal(messageString, &m); err != nil {
		log.Fatalln("FATAL: Message was not valid JSON")
	}

	log.Println(m)

	if m.sender != getUUID() {

		if err := hasValidFields(m); err != nil {
			log.Fatalln("FATAL: invalid fields in message", err)
		}

		if err := performAction(m); err != nil {
			log.Fatalln("FATAL: Cannot perform action:", err)
		}
	}
}

func getUUID() uuid {
	uuidBytes, err := ioutil.ReadFile(path.Join(pfsDir, "meta", "uuid"))
	if err != nil {
		log.Fatalln("FATAL: Cannot read UUID")
	}

	return uuid(uuidBytes)
}

func hasValidFields(m MessageData) *error {
	return nil
}

func performAction(m MessageData) *error {
	return nil
}
