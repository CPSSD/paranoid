package icserver

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
	"os"
)

// MessageChan is the channel to which incoming messages will be passed
// Attach a listener to this channel to receive messages
var MessageChan = make(chan FileSystemMessage)
var verbose = false

// FileSystemMessage is the structure which represents messages coming from the client
type FileSystemMessage struct {
	Command    string   `json:"command"`
	Args       []string `json:"args"`
	Base64Data string   `json:"data"`
	Data       []byte
}

// handleConnection accepts a connection and handles messages received through the connection
func handleConnection(conn net.Conn) {
	verboseLog("icserver new connection")
	defer verboseLog("icserver connection lost")
	for {
		buffer := make([]byte, 1024)
		mSize, err := conn.Read(buffer)
		if err != nil {
			log.Fatalln("icserver message read eror: ", err)
		}
		data := buffer[0:mSize]
		verboseLog("icserver new message:\n" + string(data))

		message := &FileSystemMessage{}
		err = json.Unmarshal(data, message)
		if err != nil {
			log.Fatalln("icserver json unmarshal error: ", err)
		}

		if len(message.Base64Data) != 0 {
			message.Data, err = base64.StdEncoding.DecodeString(message.Base64Data)
			if err != nil {
				log.Fatalln("icserver base64 decoding error: ", err)
			}
		}

		MessageChan <- *message
		verboseLog("icserver new message pushed to channel: " + message.Command)
	}
}

// RunServer runs the server
// give a true parameter for verbose logging
func RunServer(verboseLogging bool) {
	verbose = verboseLogging

	sockFIlePath := "/tmp/pfic.sock"
	listener, err := net.Listen("unix", sockFIlePath)
	if err != nil {
		log.Fatalln("ic listen error: ", err)
	}

	defer os.Remove(sockFIlePath)
	defer verboseLog("icserver no longer listening")

	verboseLog("icserver listening on " + sockFIlePath)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("ic accept error: ", err)
		}

		go handleConnection(conn)
	}
}

// verboseLog logs what the server is doing if the verboseLogging option was
// given when running the server
func verboseLog(message string) {
	if verbose {
		log.Println(message)
	}
}
