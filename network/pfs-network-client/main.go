package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
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

	url := url.URL{Scheme: "ws", Host: serverAddr, Path: "/connect"}
	log.Println("Connecting to server...")

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatalln("FATAL: Cannot Connect:", err)
	} else {
		log.Println("Establised connection with", url.String())
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("FATAL: Cannot Read Message:", err)
			continue
		}
		data, err := parseMessage(message)
		if err != nil {
			log.Println("FATAL: Cannot Parse Message:", err)
			continue
		}

		if len(data.Sender) == 0 {
			log.Println("FATAL: The Sender must be specified")
			continue
		}

		if data.Sender != GetUUID(pfsDir) {
			if err := hasValidFields(data); err != nil {
				log.Println("FATAL: invalid fields in message:", err)
				continue
			}

			if err := runPfsCommand(data, pfsDir); err != nil {
				log.Println("FATAL: Cannot perform action", data.Type, ":", err)
				continue
			}
		}
	}

	conn.Close()
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
	Data   []byte `json:"data,omitempty"`
}

// Parse Message converts the input from JSON into the MessageData struct
func parseMessage(messageString []byte) (MessageData, error) {
	m := MessageData{}

	if err := json.Unmarshal(messageString, &m); err != nil {
		log.Println(err)
		return m, err
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

// hasValidFields checks do all messages have required inputs
func hasValidFields(message MessageData) error {
	switch message.Type {
	case "write":
		if len(message.Name) == 0 {
			return errors.New("write")
		}
	case "creat":
		if len(message.Name) == 0 {
			return errors.New("creat")
		}
	case "link":
		if (len(message.Name) == 0) || (len(message.Target) == 0) {
			return errors.New("link")
		}
	case "unlink":
		if len(message.Name) == 0 {
			return errors.New("unlink")
		}
	case "truncate":
		if len(message.Name) == 0 {
			return errors.New("truncate")
		}
	default:
		return errors.New("No type provided")
	}

	return nil
}

// runPfsCommand calls pfs from the network so the changes can be performed
// accross all connected devices
func runPfsCommand(message MessageData, pfsDir string) error {
	switch message.Type {
	case "creat":
		command := exec.Command("pfs", "-n", "creat", pfsDir, message.Name)
		return command.Run() // Returns the error message
	case "write":
		command := exec.Command("pfs", "-n", "write", pfsDir, message.Name, strconv.Itoa(message.Offset), strconv.Itoa(message.Length))
		pipe, err := command.StdinPipe()
		if err != nil {
			return err
		}
		io.WriteString(pipe, string(message.Data))
		pipe.Close()
		return command.Run() // Returns the error message
	}

	return nil
}
