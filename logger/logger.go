package logger

import (
	"log"
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

func (l *paranoidLogger) Debug(v ...interface{}) {
	format := "[DEBUG] " + l.component + ":"
	args := make([]interface{}, 0)
	args = append(args, format)
	args = append(args, v...)

	log.Println(args...)
}

func (l *paranoidLogger) Debugf(format string, v ...interface{}) {
	format = "[DEBUG] " + l.component + ": " + format
	log.Printf(format, v...)
}
