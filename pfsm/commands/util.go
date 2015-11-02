package commands

import (
	"log"
	"os"
	"strconv"
)

//Current pfsm supported return codes
const (
	OK     = iota
	ENOENT //No such file or directory.
	EACCES //Can not access file
)

func getReturnCode(code int) string {
	strcode := strconv.Itoa(code)
	if len(strcode) < 2 {
		return "0" + strcode
	}
	return strcode
}

func checkFileExists(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
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
