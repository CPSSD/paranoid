package main

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/intercom"
	"github.com/cpssd/paranoid/pfsd/pnetserver"
	"github.com/cpssd/paranoid/pfsd/upnp"
	"github.com/kardianos/osext"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

func stopAllServices() {
	if globals.UPnPEnabled {
		err := upnp.ClearPortMapping(globals.Port)
		if err != nil {
			log.Info("Could not clear port mapping. Error : ", err)
		}
	}
	close(globals.Quit) // Sends stop signal to all goroutines
	if !*noNetwork {
		close(raftNetworkServer.Quit)
		srv.Stop()
		raftNetworkServer.Wait.Wait()
	}
	err := intercom.ShutdownServer()
	if err != nil {
		log.Warn("Could not shut down internal communication server:", err)
	} else {
		log.Info("Internal communication server stopped.")
	}
	// Since srv can't talk to the waitgroup itself, we do on its behalf
	// We also wait to give it some time to stop itself.
	time.Sleep(time.Millisecond * 10)
	globals.Wait.Done()
	log.Info("ParanoidNetwork server stopped.")
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
	log.Info("SIGHUP received. Restarting.")
	stopAllServices()
	log.Info("All services stopped. Forking process.")
	execSpec := &syscall.ProcAttr{
		Env: os.Environ(),
	}
	pathToSelf, err := osext.Executable()
	if err != nil {
		log.Warn("Could not get path to self:", err)
		pathToSelf = os.Args[0]
	}
	fork, err := syscall.ForkExec(pathToSelf, os.Args, execSpec)
	if err != nil {
		log.Error("Could not fork child PFSD instance:", err)
	} else {
		log.Info("Forked successfully. New PID:", fork)
	}
}

func handleSIGTERM() {
	log.Info("SIGTERM received. Exiting.")
	stopAllServices()
	err := os.Remove(path.Join(pnetserver.ParanoidDir, "meta", "pfsd.pid"))
	if err != nil {
		log.Info("Can't remove PID file ", err)
	}
	log.Info("All services stopped. Have a nice day.")
}
