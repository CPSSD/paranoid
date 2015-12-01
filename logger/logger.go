package logger

import (
	"io"
	"log"
	"os"
	"path"
)

// LogLevel is an abstraction over int that allows to better undestand the
// input of SetLogLevel
type LogLevel int

const (
	DEBUG   LogLevel = iota
	VERBOSE LogLevel = iota
	INFO    LogLevel = iota
	WARNING LogLevel = iota
	ERROR   LogLevel = iota
)

// Output enums to set the outputs
type LogOutput int

const (
	STDERR  LogOutput = iota << 1
	LOGFILE LogOutput = iota << 1
)

// Logger stuct containing the variables necessary for the logger
type paranoidLogger struct {
	component string
	curPack   string
	logDir    string
	writer    io.Writer
	logLevel  LogLevel
	native    *log.Logger
}

// New creates a new logger and returns a new logger
func New(currentPackage string, component string, logDirectory string) *paranoidLogger {
	l := paranoidLogger{
		component: component,
		curPack:   currentPackage,
		logDir:    logDirectory,
		logLevel:  INFO,
		native:    log.New(nil, "", log.LstdFlags)}

	if _, err := os.Stat(logDirectory); err != nil {
		l.Fatalf("Log directory %s not found\n", logDirectory)
	}
	l.SetOutput(STDERR)
	return &l
}

// SetLogLevel sets the logging level where the level is a constant
func (l *paranoidLogger) SetLogLevel(level LogLevel) {
	l.logLevel = level
}

// SetOutput sets the default output for the
func (l *paranoidLogger) SetOutput(output LogOutput) {
	var writers []io.Writer

	switch {
	case STDERR|LOGFILE == output:
		w, err := createFileWriter(l.logDir, l.component)
		if err != nil {
			l.Fatal("Cannot write to log file: ", err)
		}
		writers = append(writers, w)
		writers = append(writers, os.Stderr)
	case STDERR == output:
		writers = append(writers, os.Stderr)
	case LOGFILE == output:
		w, err := createFileWriter(l.logDir, l.component)
		if err != nil {
			l.Fatal("Cannot write to log file: ", err)
		}
		writers = append(writers, w)
	default:
		writers = append(writers, os.Stderr)
	}

	l.writer = io.MultiWriter(writers...)
	l.native.SetOutput(l.writer)
}

// AddAdditionalWriter allows to add a custom writer to the logger.
// This can be cleared by calling logger.SetOutput() again
func (l *paranoidLogger) AddAdditionalWriter(writer io.Writer) {
	l.writer = io.MultiWriter(l.writer, writer)
	l.native.SetOutput(l.writer)
}

///////////////////////////////// DEBUG /////////////////////////////////

// Debug only prints if DEBUG env var is set
func (l *paranoidLogger) Debug(v ...interface{}) {
	if l.logLevel <= DEBUG {
		l.output("DEBUG", v...)
	}
}

// Debug only prints if DEBUG env var is set
func (l *paranoidLogger) Debugf(format string, v ...interface{}) {
	if l.logLevel <= DEBUG {
		l.outputf("DEBUG", format, v...)
	}
}

///////////////////////////////// VERBOSE /////////////////////////////////

func (l *paranoidLogger) Verbose(v ...interface{}) {
	if l.logLevel <= VERBOSE {
		l.Info(v...)
	}
}

func (l *paranoidLogger) Verbosef(format string, v ...interface{}) {
	if l.logLevel <= VERBOSE {
		l.Infof(format, v...)
	}
}

///////////////////////////////// INFO /////////////////////////////////

// Info logs as type info
func (l *paranoidLogger) Info(v ...interface{}) {
	if l.logLevel <= INFO {
		l.output("INFO", v...)
	}
}

func (l *paranoidLogger) Infof(format string, v ...interface{}) {
	if l.logLevel <= INFO {
		l.outputf("INFO", format, v...)
	}
}

///////////////////////////////// WARN /////////////////////////////////

func (l *paranoidLogger) Warn(v ...interface{}) {
	if l.logLevel <= WARNING {
		l.output("WARN", v...)
	}
}

func (l *paranoidLogger) Warnf(format string, v ...interface{}) {
	if l.logLevel <= WARNING {
		l.outputf("WARN", format, v...)
	}
}

///////////////////////////////// ERROR /////////////////////////////////

func (l *paranoidLogger) Error(v ...interface{}) {
	if l.logLevel <= ERROR {
		l.output("ERROR", v...)
	}
}

func (l *paranoidLogger) Errorf(format string, v ...interface{}) {
	if l.logLevel <= ERROR {
		l.outputf("ERROR", format, v...)
	}
}

///////////////////////////////// FATAL /////////////////////////////////

func (l *paranoidLogger) Fatal(v ...interface{}) {
	l.output("FATAL", v...)
	os.Exit(1)
}

func (l *paranoidLogger) Fatalf(format string, v ...interface{}) {
	l.outputf("FATAL", format, v...)
	os.Exit(1)
}

///////////////////////////////// GENERAL /////////////////////////////////

func (l *paranoidLogger) output(mtype string, v ...interface{}) {
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

func (l *paranoidLogger) outputf(mtype string, format string, v ...interface{}) {
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
