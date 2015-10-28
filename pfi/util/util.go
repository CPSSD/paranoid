package util

import (
	"log"
)

var LogOutput bool
var MountPoint string
var PfsInitPoint string

//LogMessage checks if the -v flag was specified and either logs or doesnt log the message
func LogMessage(message string) {
	if LogOutput {
		log.Println(message)
	}
}
