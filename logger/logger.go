package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

// LogLevel is an abstraction over int that allows to better undestand the
// input of SetLogLevel
type LogLevel int

const (
	DEBUG LogLevel = iota
	VERBOSE
	INFO
	WARNING
	ERROR
)

// Output enums to set the outputs
type LogOutput int

const (
	STDERR  LogOutput = 1 << (iota + 1)
	LOGFILE LogOutput = 1 << (iota + 1)
)

// ParanoidLogger struct containing the variables necessary for the logger
type ParanoidLogger struct {
	component string
	curPack   string
	logDir    string
	writer    io.Writer
	logLevel  LogLevel
	native    *log.Logger
}

// New creates a new logger and returns a new logger
func New(currentPackage string, component string, logDirectory string) *ParanoidLogger {
	l := ParanoidLogger{
		component: component,
		curPack:   currentPackage,
		logDir:    logDirectory,
		writer:    os.Stderr,
		logLevel:  INFO,
		native:    log.New(os.Stderr, "", log.LstdFlags),
	}
	return &l
}

// SetLogLevel sets the logging level for the logger
func (l *ParanoidLogger) SetLogLevel(level LogLevel) {
	l.logLevel = level
}

// SetOutput sets the default output for the logger
func (l *ParanoidLogger) SetOutput(output LogOutput) error {
	var writers []io.Writer

	if STDERR&output == STDERR {
		writers = append(writers, os.Stderr)
	}
	if LOGFILE&output == LOGFILE {
		w, err := createFileWriter(l.logDir, l.component)
		if err != nil {
			return fmt.Errorf("failed accessing log file: %s", err)
		}
		writers = append(writers, w)
	}

	l.writer = io.MultiWriter(writers...)
	l.native.SetOutput(l.writer)
	return nil
}

// AddAdditionalWriter allows to add a custom writer to the logger.
// This can be cleared by calling logger.SetOutput() again
func (l *ParanoidLogger) AddAdditionalWriter(writer io.Writer) {
	l.writer = io.MultiWriter(l.writer, writer)
	l.native.SetOutput(l.writer)
}

///////////////////////////////// DEBUG /////////////////////////////////

// Debug only prints if LogLevel is set to DEBUG
func (l *ParanoidLogger) Debug(v ...interface{}) {
	if l.logLevel <= DEBUG {
		l.output("DEBUG", v...)
	}
}

// Debug only prints if LogLevel is set to DEBUG
func (l *ParanoidLogger) Debugf(format string, v ...interface{}) {
	if l.logLevel <= DEBUG {
		l.outputf("DEBUG", format, v...)
	}
}

///////////////////////////////// VERBOSE /////////////////////////////////

// Verbose only prints if LogLevel is set to VERBOSE or lower in importance
func (l *ParanoidLogger) Verbose(v ...interface{}) {
	if l.logLevel <= VERBOSE {
		l.Info(v...)
	}
}

// Verbose only prints if LogLevel is set to VERBOSE or lower in importance
func (l *ParanoidLogger) Verbosef(format string, v ...interface{}) {
	if l.logLevel <= VERBOSE {
		l.Infof(format, v...)
	}
}

///////////////////////////////// INFO /////////////////////////////////

// Info only prints if LogLevel is set to INFO or lower in importance
func (l *ParanoidLogger) Info(v ...interface{}) {
	if l.logLevel <= INFO {
		l.output("INFO", v...)
	}
}

// Info only prints if LogLevel is set to INFO or lower in importance
func (l *ParanoidLogger) Infof(format string, v ...interface{}) {
	if l.logLevel <= INFO {
		l.outputf("INFO", format, v...)
	}
}

///////////////////////////////// WARN /////////////////////////////////

// Warn only prints if LogLevel is set to WARNING or lower in importance
func (l *ParanoidLogger) Warn(v ...interface{}) {
	if l.logLevel <= WARNING {
		l.output("WARN", v...)
	}
}

// Warn only prints if LogLevel is set to WARNING or lower in importance
func (l *ParanoidLogger) Warnf(format string, v ...interface{}) {
	if l.logLevel <= WARNING {
		l.outputf("WARN", format, v...)
	}
}

///////////////////////////////// ERROR /////////////////////////////////

// Error only prints if LogLevel is set to ERROR or lower in importance
func (l *ParanoidLogger) Error(v ...interface{}) {
	if l.logLevel <= ERROR {
		l.output("ERROR", v...)
	}
}

// Error only prints if LogLevel is set to ERROR or lower in importance
func (l *ParanoidLogger) Errorf(format string, v ...interface{}) {
	if l.logLevel <= ERROR {
		l.outputf("ERROR", format, v...)
	}
}

///////////////////////////////// FATAL /////////////////////////////////

// Fatal always prints and exits the program with exit code 1
func (l *ParanoidLogger) Fatal(v ...interface{}) {
	l.output("FATAL", v...)
	os.Exit(1)
}

// Fatal always prints and exits the program with exit code 1
func (l *ParanoidLogger) Fatalf(format string, v ...interface{}) {
	l.outputf("FATAL", format, v...)
	os.Exit(1)
}

///////////////////////////////// GENERAL /////////////////////////////////

func (l *ParanoidLogger) output(mtype string, v ...interface{}) {
	fmt := "[" + mtype + "] "
	// Add an extra space if the message type (mtype) is only 4 letters long
	if len(mtype) == 4 {
		fmt += " " + l.curPack + ":"
	} else {
		fmt += l.curPack + ":"
	}

	var args []interface{}
	args = append(args, fmt)
	args = append(args, v...)

	l.native.Println(args...)
}

func (l *ParanoidLogger) outputf(mtype string, format string, v ...interface{}) {
	fmt := "[" + mtype + "] "
	// Add an extra space if the message type (mtype) is only 4 letters long
	if len(mtype) == 4 {
		fmt += " " + l.curPack + ": " + format
	} else {
		fmt += l.curPack + ": " + format
	}

	l.native.Printf(fmt, v...)
}

func createFileWriter(logPath string, component string) (io.Writer, error) {
	return os.OpenFile(path.Join(logPath, component+".log"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
}
