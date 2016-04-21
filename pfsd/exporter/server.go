package exporter

import (
	"golang.org/x/net/websocket"
	"net/http"
)

type Server struct {
	messages chan Message
	port     string
	client   *client
}

func NewServer(port string) *Server {
	return &Server{
		port:     port,
		messages: make(chan Message),
	}
}

// Run starts the websocket server
func (s *Server) Run() {
	onConnected := func(ws *websocket.Conn) {
		s.client = &client{
			ws:   ws,
			msgs: make(chan Message),
		}
		s.client.listen()

		// Send the state message
		s.Send(Message{
			Type: StateMessage,
			Data: MessageData{
				Nodes: toNodeArray(nodeList),
			},
		})
	}

	http.ListenAndServe(":"+s.port, websocket.Handler(onConnected))
}

func (s *Server) Send(msg Message) {
	s.client.write(msg)
}
