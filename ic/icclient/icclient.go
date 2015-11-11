package icclient

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
	"path"
)

// fileSystemMessage is the container for messages to be sent to server
type fileSystemMessage struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Data    string   `json:"data"`
}

// SendMessage sends a message to the server
func SendMessage(pfsDirectory string, command string, arguments []string) {
	message := fileSystemMessage{
		Command: command,
		Args:    arguments,
		Data:    "",
	}

	dialAndSend(pfsDirectory, message)
}

// SendMessageWithData sends a message with data to the server
func SendMessageWithData(pfsDirectory string, command string, arguments []string, data []byte) {
	base64Data := base64.StdEncoding.EncodeToString(data)

	message := fileSystemMessage{
		Command: command,
		Args:    arguments,
		Data:    base64Data,
	}

	dialAndSend(pfsDirectory, message)
}

// dialAndSend dials the server and sends a message
func dialAndSend(pfsDirectory string, message fileSystemMessage) {
	sockFilePath := path.Join(pfsDirectory, "meta", "pfic.sock")

	messageData, err := json.Marshal(message)
	if err != nil {
		log.Fatalln("icclient json Marshal error: ", err)
	}

	conn, err := net.Dial("unix", sockFilePath)
	if err != nil {
		log.Fatalln("icclient connection error: ", err)
	}

	_, err = conn.Write(messageData)

	err = conn.Close()
	if err != nil {
		log.Fatalln("icclient socket close error: ", err)
	}
}
