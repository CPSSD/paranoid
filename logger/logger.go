package logger

import (
	"log"
	"os"
)

type loggerFlags struct {
	verbose bool
	debug   bool
	output  string
}

// Logger stuct containing the variables necessary for the logger
type paranoidLogger struct {
	component string
	curPack   string
	logDir    string
	flags     loggerFlags
}

// New creates a new logger and returns a new logger
func New(component string, currentPackage string, logDirectory string) paranoidLogger {
	return paranoidLogger{
		component: component,
		curPack:   currentPackage,
		logDir:    logDirectory,
		flags: loggerFlags{
			debug:  os.Getenv("DEBUG") == "true",
			output: "both"}}
}

func (l *paranoidLogger) SetFlag(flag string, value bool) bool {
	switch flag {
	case "verbose":
		l.flags.verbose = value
	case "debug":
		l.flags.debug = value
	default:
		return false
	}
	return true
}

// Info logs as type info
func (l *paranoidLogger) Info(v ...interface{}) {
	format := "[INFO]  " + l.component + ":"
	args := make([]interface{}, 0)
	args = append(args, format)
	args = append(args, v...)

	log.Println(args...)
}

func (l *paranoidLogger) Infof(format string, v ...interface{}) {
	format = "[INFO]  " + l.component + ": " + format
	log.Printf(format, v...)
}

func (l *paranoidLogger) Warn(v ...interface{}) {
	format := "[WARN]  " + l.component + ":"
	args := make([]interface{}, 0)
	args = append(args, format)
	args = append(args, v...)

	log.Println(args...)
}

func (l *paranoidLogger) Warnf(format string, v ...interface{}) {
	format = "[WARN]  " + l.component + ": " + format
	log.Printf(format, v...)
}

func (l *paranoidLogger) Error(v ...interface{}) {
	format := "[ERROR] " + l.component + ":"
	args := make([]interface{}, 0)
	args = append(args, format)
	args = append(args, v...)

	log.Println(args...)
}

func (l *paranoidLogger) Errorf(format string, v ...interface{}) {
	format = "[ERROR] " + l.component + ": " + format
	log.Printf(format, v...)
}

// Debug only prints if DEBUG env var is set
func (l *paranoidLogger) Debug(v ...interface{}) {
	if !l.flags.debug {
		return
	}
	format := "[DEBUG] " + l.component + ":"
	args := make([]interface{}, 0)
	args = append(args, format)
	args = append(args, v...)

	log.Println(args...)
}

// Debug only prints if DEBUG env var is set
func (l *paranoidLogger) Debugf(format string, v ...interface{}) {
	if !l.flags.debug {
		return
	}
	format = "[DEBUG] " + l.component + ": " + format
	log.Printf(format, v...)
}

func (l *paranoidLogger) Fatal(v ...interface{}) {
	format := "[FATAL] " + l.component + ":"
	args := make([]interface{}, 0)
	args = append(args, format)
	args = append(args, v...)

	log.Fatalln(args...)
}

func (l *paranoidLogger) Fatalf(format string, v ...interface{}) {
	format = "[FATAL] " + l.component + ": " + format
	log.Fatalf(format, v...)
}

func (l *paranoidLogger) Verbose(v ...interface{}) {
	if l.flags.verbose {
		l.Info(v...)
	}
}

func (l *paranoidLogger) Verbosef(format string, v ...interface{}) {
	if l.flags.verbose {
		l.Infof(format, v...)
	}
}
