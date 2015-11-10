package icserver

import (
	"log"
	"net"
	"os"
)

// RunServer runs the server
func RunServer() {
	sockFIlePath := "/tmp/pfic.sock"
	listener, err := net.Listen("unix", sockFIlePath)
	if err != nil {
		log.Fatalln("ic listen error: ", err)
	}
	defer os.Remove(sockFIlePath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("ic accept error: ", err)
		}

		// send somewhere
	}
}
