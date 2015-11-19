package logger

import (
	"log"
	"os"
)

// Logger stuct containing the variables necessary for the logger
type paranoidLogger struct {
	component string
	curPack   string
	logDir    string
}

// New creates a new logger and returns a new logger
func New(component string, currentPackage string, logDirectory string) paranoidLogger {
	return paranoidLogger{component: component, curPack: currentPackage, logDir: logDirectory}
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
	if os.Getenv("DEBUG") == "true" {
		format := "[DEBUG] " + l.component + ":"
		args := make([]interface{}, 0)
		args = append(args, format)
		args = append(args, v...)

		log.Println(args...)
	}

}

// Debug only prints if DEBUG env var is set
func (l *paranoidLogger) Debugf(format string, v ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		format = "[DEBUG] " + l.component + ": " + format
		log.Printf(format, v...)
	}
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
