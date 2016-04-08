package main

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/upnp"
	"github.com/kardianos/osext"
	"os"
	"os/signal"
	"path"
	"syscall"
)

func stopAllServices() {
	globals.ShuttingDown = true
	if globals.UPnPEnabled {
		err := upnp.ClearPortMapping(globals.ThisNode.Port)
		if err != nil {
			log.Info("Could not clear port mapping. Error : ", err)
		}
	}
	close(globals.Quit) // Sends stop signal to all goroutines

	// Save all KeyPieces to disk, to ensure we haven't missed any so far.
	globals.HeldKeyPieces.SaveToDisk()

	if !globals.NetworkOff {
		close(globals.RaftNetworkServer.Quit)
		srv.Stop()
		globals.RaftNetworkServer.Wait.Wait()
	}

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
	err := os.Remove(path.Join(globals.ParanoidDir, "meta", "pfsd.pid"))
	if err != nil {
		log.Info("Can't remove PID file ", err)
	}
	log.Info("All services stopped. Have a nice day.")
}
