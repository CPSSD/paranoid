package main

import (
	"flag"
	"log"
	"sync"
)

var (
	portFlag = flag.String("port", "10100", "port on which to run the sample test server")
	demoFlag = flag.Bool("demo", false, "start the demo")
)

var (
	server   *Server
	thisNode Node
	wait     sync.WaitGroup
)

func main() {
	flag.Parse()

	// Initialize this node
	thisNode = Node{
		CommonName: "node1",
		Addr:       "10.0.0.1:77777",
		State:      Current,
		UUID:       "1234-abcd-5678-efgh",
	}

	wait.Add(2)
	server = NewServer(*portFlag)
	go server.Listen()

	log.Println("Waiting...")
	wait.Wait()
	log.Println("Done Waiting")
	wait.Add(1)

	if *demoFlag {
		go demo(server)
	}

	for {

	}
}
