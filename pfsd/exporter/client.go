package exporter

import (
	"golang.org/x/net/websocket"
)

type client struct {
	ws   *websocket.Conn
	msgs chan Message
}

func (c *client) write(msg Message) {
	c.msgs <- msg
}

func (c *client) listen() {
	for {
		select {
		case msg := <-c.msgs:
			websocket.JSON.Send(c.ws, msg)
		}
	}
}
