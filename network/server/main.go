package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

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

func (conn *connection) Reader() {
	for {
		_, message, err := conn.ws.ReadMessage()
		if err != nil {
			log.Println("ERROR: Could not read message: ", err)
			break
		}
		conn.p.messages <- message
	}
	conn.ws.Close()
}

func (conn *connection) Writer() {
	for message := range conn.messages {
		err := conn.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("ERROR: Could not write to websocket: ", err)
			break
		}
	}
	conn.ws.Close()
}

type websocketHandler struct {
	p *connPool
}

func (handler websocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ERROR: Could not upgrade websocket: ", err)
		return
	}
	conn := NewConnection(ws, handler.p)
	handler.p.Register <- conn
	defer func() { handler.p.Unregister <- conn }()
	go conn.Writer()
	conn.Reader()
}

func main() {
	pool := NewConnPool()
	go pool.Run()
	http.Handle("/connect", websocketHandler{p: pool})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("FATAL: ", err)
	}
}
