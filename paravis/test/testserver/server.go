package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

// Server contains the simple server for the
type Server struct {
	Port     string
	conn     *websocket.Conn
	messages chan Message
}

// NewServer creates a pointer to a new server
func NewServer(port string) *Server {
	return &Server{
		Port:     port,
		messages: make(chan Message, 10),
	}
}

// Listen starts the server on the specified port
func (s *Server) Listen() {
	wait.Done()
	err := http.ListenAndServe(":"+s.Port, websocket.Handler(s.onConnection))
	if err != nil {
		log.Print("Error:", err)
	}
}

// Send just sends to the servers client
func (s *Server) Send(m Message) {
	s.messages <- m
}

// onConnection handles a connection from the client
func (s *Server) onConnection(conn *websocket.Conn) {
	log.Println("Client connecting...")
	s.conn = conn
	wait.Done()
	s.Send(Message{
		Type: TypeState,
		Data: Data{
			Nodes: []Node{thisNode},
		},
	})
	for {
		select {
		case m := <-s.messages:
			log.Println("Sending message...")
			err := websocket.JSON.Send(s.conn, m)
			if err != nil {
				log.Println("ERROR: Client write:", err)
			}
		}
	}
}
