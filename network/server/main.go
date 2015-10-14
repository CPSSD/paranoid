package main

import (
	"github.com/gorilla/websocket"
)

type connPool struct {
	connections map[*connection]bool

	// Buffered channel of incoming messages
	messages chan []byte

	// Incoming connections
	Register chan *connection

	// Requests to deactivate a connection
	Unregister chan *connection
}

func NewConnPool() *connPool {
	return &connPool{
		connections: make(map[*connection]bool),
		messages:    make(chan []byte, 256),
		Register:    make(chan *connection),
		Unregister:  make(chan *connection),
	}
}

func (p *connPool) Run() {
	for {
		select {
		case conn := <-p.Register:
			p.connections[conn] = true
		case conn := <-p.Unregister:
			if _, ok := p.connections[conn]; ok {
				delete(p.connections, conn)
			}
		// When receiving a message, send it down all active connections
		case message := <-p.messages:
			for conn := range p.connections {
				conn.messages <- message
			}
		}
	}
}

type connection struct {
	ws *websocket.Conn
	p  *connPool
	// Buffered channel of messages to be send to the client
	messages chan []byte
}

func NewConnection(ws *websocket.Conn, p *connPool) *connection {
	return &connection{
		ws:       ws,
		p:        p,
		messages: make(chan []byte, 256),
	}
}
