package main

import (
	"os"
	"os/signal"
	"syscall"
)

// Listens for SIGTERM and SIGHUP. Should be run in own goroutine.
func HandleSignals() {
	incoming := make(chan os.Signal, 1)
	signal.Notify(incoming, syscall.SIGHUP)
	for {
		sig := <-incoming
		switch sig {
		case syscall.SIGHUP:
			handleSIGHUP()
		}
	}
}

func handleSIGHUP() {
	return
}
