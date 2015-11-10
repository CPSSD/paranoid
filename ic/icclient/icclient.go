package icclient

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
)

// fileSystemMessage is the container for messages to be sent to server
type fileSystemMessage struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Data    string   `json:"data"`
}

// SendMessage sends a message to the server
func SendMessage(command string, arguments []string) {
	message := &fileSystemMessage{
		Command: command,
		Args:    arguments,
		Data:    "",
	}

	dialAndSend(*message)
}

// SendMessageWithData sends a message with data to the server
func SendMessageWithData(command string, arguments []string, data []byte) {
	base64Data := base64.StdEncoding.EncodeToString(data)

	message := &fileSystemMessage{
		Command: command,
		Args:    arguments,
		Data:    base64Data,
	}

	dialAndSend(*message)
}

// dialAndSend dials the server and sends a message
func dialAndSend(message fileSystemMessage) {
	messageData, err := json.Marshal(message)
	if err != nil {
		log.Fatalln("icclient json Marshal error: ", err)
	}

	conn, err := net.Dial("unix", "/tmp/pfic.sock")
	if err != nil {
		log.Fatalln("icclient connection error: ", err)
	}

	_, err = conn.Write(messageData)

	err = conn.Close()
	if err != nil {
		log.Fatalln("icclient socket close error: ", err)
	}
}
