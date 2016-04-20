// exporter package creates an WebServer which exports Raft information
package exporter

import (
	"github.com/cpssd/paranoid/logger"
)

var Log *logger.ParanoidLogger
var server *Server

func Send(msg Message) {
	server.Send(msg)
}

func NewStdServer(port string) {
	server = NewServer(port)
}

func Listen() {
	server.Run()
}
