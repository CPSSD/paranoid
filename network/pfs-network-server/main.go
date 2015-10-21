package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strconv"
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

func (pool *connPool) Run() {
	for {
		select {
		case conn := <-pool.Register:
			pool.connections[conn] = true
		case conn := <-pool.Unregister:
			if _, ok := pool.connections[conn]; ok {
				delete(pool.connections, conn)
			}
		// When receiving a message, send it down all active connections
		case message := <-pool.messages:
			log.Println("DEBUG: message received:", message)
			for conn := range pool.connections {
				conn.messages <- message
			}
		}
	}
}

type connection struct {
	websocket *websocket.Conn
	pool      *connPool
	// Buffered channel of messages to be send to the client
	messages chan []byte
}

func NewConnection(websocket *websocket.Conn, pool *connPool) *connection {
	return &connection{
		websocket: websocket,
		pool:      pool,
		messages:  make(chan []byte, 256),
	}
}

func (conn *connection) Reader() {
	for {
		_, message, err := conn.websocket.ReadMessage()
		if err != nil {
			log.Println("ERROR: Could not read message: ", err)
			break
		}
		conn.pool.messages <- message
	}
	conn.websocket.Close()
}

func (conn *connection) Writer() {
	for message := range conn.messages {
		err := conn.websocket.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("ERROR: Could not write to websocket: ", err)
			break
		}
	}
	conn.websocket.Close()
}

type websocketHandler struct {
	pool *connPool
}

func (handler websocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	websocket, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ERROR: Could not upgrade websocket: ", err)
		return
	}
	conn := NewConnection(websocket, handler.pool)
	handler.pool.Register <- conn
	defer func() { handler.pool.Unregister <- conn }()
	go conn.Writer()
	conn.Reader()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print("Usage:\n\tpfs-network-server <port>\n")
		os.Exit(1)
	}
	port := os.Args[1]
	portint, err := strconv.Atoi(port)
	if err != nil || portint < 1 || portint > 65535 {
		log.Fatalln("FATAL: port must be a number between 1 and 65535, inclusive.")
	}

	pool := NewConnPool()
	go pool.Run()
	http.Handle("/connect", websocketHandler{pool: pool})
	log.Println("INFO: Server listening on port ", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalln("FATAL: ", err)
	}
}
