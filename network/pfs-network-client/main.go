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
	"os/exec"
	"path"
	"strconv"
)

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
		data, err := ParseMessage(message)
		if err != nil {
			log.Fatalln(err)
		}

		if len(data.Sender) == 0 {
			log.Fatalln("FATAL: The Sender must be specified")
		}

		if data.Sender != GetUUID(pfsDir) {
			if err := HasValidFields(data); err != nil {
				log.Fatalln("FATAL: invalid fields in message:", err)
			}

			if err := PerformAction(data, pfsDir); err != nil {
				log.Fatalln("FATAL: Cannot perform action", data.Type, ":", err)
			}
		}
	}

	defer c.Close()
}

// MessageData stores all the values that the server can provide it.
// It ommits values if something is empty
type MessageData struct {
	Sender string `json:"sender,omitempty"`
	Type   string `json:"type,omitempty"`
	Name   string `json:"name,omitempty"`
	Target string `json:"target,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Length int    `json:"length,omitempty"`
	Data   string `json:"data,omitempty"`
}

// Parse Message converts the input from JSON into the MessageData struct
func ParseMessage(messageString []byte) (MessageData, error) {
	m := MessageData{}

	if err := json.Unmarshal(messageString, &m); err != nil {
		log.Println(err)
		return m, errors.New("FATAL: Message was not valid JSON")
	}

	return m, nil
}

// GetUUID takes the uuid of the pfs from the meta file
func GetUUID(pfsDir string) string {
	uuidBytes, err := ioutil.ReadFile(path.Join(pfsDir, "meta", "uuid"))
	if err != nil {
		log.Fatalln("FATAL: Cannot read UUID:", err)
	}

	return string(uuidBytes)
}

// HasValidFields checks do all messages have required inputs
func HasValidFields(m MessageData) error {
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

// PerformAction calls pfs from the network so the changes can be performed
// accross all connected devices
func PerformAction(m MessageData, pfsDir string) error {
	switch m.Type {
	case "creat":
		command := exec.Command("pfs", "-n", "creat", pfsDir, m.Name)
		if err := command.Run(); err != nil {
			return err
		}
	//case "write":
		// TODO: Implement write with pipes

		//command := exec.Command("pfs", "-n", "write", pfsDir, m.Name, m.Offset, m.Length)
	}

	return nil
}
