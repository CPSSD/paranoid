package network

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

func sendMessage(message []byte, host, port string) {
	url := url.URL{Scheme: "ws", Host: host + ":" + port, Path: "/connect"}
	log.Printf("connecting to %s", url.String())

	ws, _, err := websocket.DefaultDialer.Dial(url.String(), nil)

	if err != nil {
		log.Fatal("dial:", err)
	}

	defer ws.Close()

	err = ws.WriteJSON(message)
	if err != nil {
		log.Fatal(err)
	}
}
