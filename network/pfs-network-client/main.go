package main

import (
	//"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"strconv"
)

type messageData struct {
	Sender string `json:"sender,omitempty"`
	Type   string `json:"type,omitempty"`
	Name   string `json:"name,omitempty"`
	Target string `json:"target,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Length int    `json:"length,omitempty"`
	Data   string `json:"data,omitempty"`
}

func main() {
	args := os.Args[1:]

	if len(args) != 4 {
		fmt.Println("Usage:\n\tpfs-network-client --client <pfs-directory> <server-ip> <server-port>")
		os.Exit(1)
	}

	pfsDir := args[1]

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
		data, err := parseMessage(message)
		if err != nil {
			log.Fatalln(err)
		}

		if len(data.Sender) == 0 {
			log.Fatalln("FATAL: The Sender must be specified")
		}

		if data.Sender != getUUID(pfsDir) {
			if err := hasValidFields(data); err != nil {
				log.Fatalln("FATAL: invalid fields in message:", err)
			}

			if err := performAction(data); err != nil {
				log.Fatalln("FATAL: Cannot perform action", data.Type)
			}
		}
	}

	defer c.Close()
}

func parseMessage(messageString []byte) (messageData, error) {
	m := messageData{}

	if err := json.Unmarshal(messageString, &m); err != nil {
		log.Println(err)
		return m, errors.New("FATAL: Message was not valid JSON")
	}

	return m, nil
}

func getUUID(pfsDir string) string {
	uuidBytes, err := ioutil.ReadFile(path.Join(pfsDir, "meta", "uuid"))
	if err != nil {
		log.Fatalln("FATAL: Cannot read UUID:", err)
	}

	return string(uuidBytes)
}

func hasValidFields(m messageData) error {
	switch m.Type {
	case "write":
		if (len(m.Name) == 0) || (len(m.Data) == 0) {
			return errors.New("write")
		}
	case "creat":
		if len(m.Name) == 0 {
			return errors.New("creat")
		}
	case "link":
		if (len(m.Name) == 0) || (len(m.Target) == 0) {
			return errors.New("link")
		}
	case "unlink":
		if len(m.Name) == 0 {
			return errors.New("unlink")
		}
	case "truncate":
		if len(m.Name) == 0 {
			return errors.New("truncate")
		}
	}

	return nil
}

func performAction(m messageData) *error {
	return nil
}
