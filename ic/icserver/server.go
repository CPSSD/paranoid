package icserver

import (
	"log"
	"net"
	"os"
)

var verbose = false

// handleConnection accepts a connection and handles messages received through the connection
func handleConnection(conn net.Conn) {
	verboseLog("icserver new connection")
	defer verboseLog("icserver connection lost")
	for {

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
