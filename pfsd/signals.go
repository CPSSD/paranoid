package main

import (
	"github.com/cpssd/paranoid/pfsd/dnetclient"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/icserver"
	"github.com/cpssd/paranoid/pfsd/pnetserver"
	"github.com/cpssd/paranoid/pfsd/upnp"
	"github.com/kardianos/osext"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

func stopAllServices() {
	upnp.ClearPortMapping(globals.Port)
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
}

// HandleSignals listens for SIGTERM and SIGHUP, and dispatches to handler
// functions when a signal is received.
func HandleSignals() {
	incoming := make(chan os.Signal, 1)
	signal.Notify(incoming, syscall.SIGHUP, syscall.SIGTERM)
	sig := <-incoming
	switch sig {
	case syscall.SIGHUP:
		handleSIGHUP()
	case syscall.SIGTERM:
		handleSIGTERM()
	}
}

func handleSIGHUP() {
	log.Println("INFO: SIGHUP received. Restarting.")
	stopAllServices()
	log.Println("INFO: All services stopped. Forking process.")
	execSpec := &syscall.ProcAttr{
		Env: os.Environ(),
	}
	pathToSelf, err := osext.Executable()
	if err != nil {
		log.Println("WARNING: Could not get path to self:", err)
		pathToSelf = os.Args[0]
	}
	fork, err := syscall.ForkExec(pathToSelf, os.Args, execSpec)
	if err != nil {
		log.Println("ERROR: Could not fork child PFSD instance:", err)
	} else {
		log.Println("INFO: Forked successfully. New PID:", fork)
	}
}

func handleSIGTERM() {
	log.Println("INFO: SIGTERM received. Exiting.")
	stopAllServices()
	err := os.Remove(path.Join(pnetserver.ParanoidDir, "meta", "pfsd.pid"))
	if err != nil {
		log.Println("INFO: Can't remove PID file ", err)
	}
	log.Println("INFO: All services stopped. Have a nice day.")
}
