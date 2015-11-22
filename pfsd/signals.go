package main

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Listens for SIGTERM and SIGHUP. Should be run in own goroutine.
func HandleSignals() {
	defer globals.Wait.Done()
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
	close(globals.Quit) // Sends stop signal to all goroutines
	srv.Stop()
	// Since srv can't talk to the waitgroup itself, we do on its behalf
	// We also wait to give it some time to stop itself.
	time.Sleep(time.Millisecond * 10)
	globals.Wait.Done()
	log.Println("INFO: gRPC server stopped. Forking new process.")
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
