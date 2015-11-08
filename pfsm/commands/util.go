package commands

import (
	"log"
	"os"
	"syscall"
)

//Check if a given file exists
func checkFileExists(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

func getAccessMode(flags uint32) uint32 {
	if flags == syscall.O_RDONLY {
		return 4
	} else if flags == syscall.O_WRONLY {
		return 2
	} else if flags == syscall.O_RDWR {
		return 6
	}
	return 7
}

//verboseLog logs a message if the verbose command line flag was set.
func verboseLog(message string) {
	if Flags.Verbose {
		log.Println(message)
	}
}

//checkErr stops the execution of the program if the given error is not nil.
//Specifies the command where the error occured as cmd
func checkErr(cmd string, err error) {
	if err != nil {
		log.Fatalln(cmd, " error occured: ", err)
	}
}
