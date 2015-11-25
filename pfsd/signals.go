package main

import (
	"github.com/cpssd/paranoid/pfsd/dnetclient"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/icserver"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Listens for SIGTERM and SIGHUP. Should be run in own goroutine.
func HandleSignals() {
	incoming := make(chan os.Signal, 1)
	signal.Notify(incoming, syscall.SIGHUP)
	sig := <-incoming
	switch sig {
	case syscall.SIGHUP:
		handleSIGHUP()
	}
}

func handleSIGHUP() {
	log.Println("INFO: SIGHUP received. Restarting.")
	close(globals.Quit)     // Sends stop signal to all goroutines
	dnetclient.Disconnect() // Disconnect from the discovery server
	icserver.StopAccept()
	srv.Stop()
	// Since srv can't talk to the waitgroup itself, we do on its behalf
	// We also wait to give it some time to stop itself.
	time.Sleep(time.Millisecond * 10)
	globals.Wait.Done()
	log.Println("INFO: ParanoidNetwork server stopped.")
	globals.Wait.Wait()
	log.Println("INFO: All services stopped. Forking process.")
	execSpec := &syscall.ProcAttr{
		Env: os.Environ(),
	}
	fork, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
	if err != nil {
		log.Println("ERROR: Could not fork child PFSD instance:", err)
	} else {
		log.Println("INFO: Forked successfully. New PID:", fork)
	}
}
