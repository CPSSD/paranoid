package main

import (
	"log"
	"time"
)

func demo(server *Server) {
	log.Println("Starting demo...")
	var m Message
	time.Sleep(100 * time.Microsecond)

	// Add a new node
	m = Message{
		Type: TypeNodeChange,
		Data: Data{
			Action: "add",
			Node: Node{
				CommonName: "node2",
				Addr:       "10.0.0.2:7777",
				State:      "leader",
				UUID:       "4567-efgh-8901-ijkl",
			},
		},
	}
	log.Println("Sending first message")
	server.Send(m)
	time.Sleep(2 * time.Second)

	// Start ticking and finish after 20 times
	tick := time.NewTicker(500 * time.Millisecond)
	for i := 0; i < 20; i++ {
		select {
		case <-tick.C:
			// Send a request write
			m = Message{
				Type: TypeEvent,
				Data: Data{
					Event: Event{
						Source:  "1234-abcd-5678-efgh",
						Target:  "4567-efgh-8901-ijkl",
						Details: "write-request",
					},
				},
			}
			server.Send(m)
			log.Print("Sending write-request")
			time.Sleep(100 * time.Second)

			// Send a write
			m = Message{
				Type: TypeEvent,
				Data: Data{
					Event: Event{
						Source:  "4567-efgh-8901-ijkl",
						Target:  "1234-abcd-5678-efgh",
						Details: "write",
					},
				},
			}
			server.Send(m)
			log.Print("Sending write")
		}
	}
	tick.Stop()
}
