package logger

import (
	"io"
	"log"
	"os"
	"path"
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
	writer    io.Writer
}

// New creates a new logger and returns a new logger
func New(component string, currentPackage string, logDirectory string) *paranoidLogger {
	l := paranoidLogger{
		component: component,
		curPack:   currentPackage,
		logDir:    logDirectory,
		flags: loggerFlags{
			debug:  os.Getenv("DEBUG") == "true",
			output: "stderr"}}

	if _, err := os.Stat(logDirectory); err != nil {
		l.Fatalf("Log directory %s not found\n", logDirectory)
	}
	l.SetOutput(l.flags.output)
	return &l
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

// SetOutput sets the default output for the
func (l *paranoidLogger) SetOutput(output string) {
	l.flags.output = output
	var writers []io.Writer

	switch l.flags.output {
	case "both":
		w, err := createFileWriter(l.logDir, l.curPack)
		if err != nil {
			l.Fatal("Cannot write to log file: ", err)
		}
		writers = append(writers, w)
		writers = append(writers, os.Stderr)
	case "stderr":
		writers = append(writers, os.Stderr)
	case "logfile":
		w, err := createFileWriter(l.logDir, l.curPack)
		if err != nil {
			l.Fatal("Cannot write to log file: ", err)
		}
		writers = append(writers, w)
	default:
		writers = append(writers, os.Stderr)
	}

	l.writer = io.MultiWriter(writers...)
	log.SetOutput(l.writer)
}

// Info logs as type info
func (l *paranoidLogger) Info(v ...interface{}) {
	l.output("INFO", v...)
}

func (l *paranoidLogger) Infof(format string, v ...interface{}) {
	l.outputf("INFO", format, v...)
}

func (l *paranoidLogger) Warn(v ...interface{}) {
	l.output("WARN", v...)
}

func (l *paranoidLogger) Warnf(format string, v ...interface{}) {
	l.outputf("WARN", format, v...)
}

func (l *paranoidLogger) Error(v ...interface{}) {
	l.output("ERROR", v...)
}

func (l *paranoidLogger) Errorf(format string, v ...interface{}) {
	l.outputf("ERROR", format, v...)
}

// Debug only prints if DEBUG env var is set
func (l *paranoidLogger) Debug(v ...interface{}) {
	if !l.flags.debug {
		return
	}
	l.output("DEBUG", v...)
}

// Debug only prints if DEBUG env var is set
func (l *paranoidLogger) Debugf(format string, v ...interface{}) {
	if !l.flags.debug {
		return
	}
	l.outputf("DEBUG", format, v...)
}

func (l *paranoidLogger) Fatal(v ...interface{}) {
	l.output("FATAL", v...)
	os.Exit(1)
}

func (l *paranoidLogger) Fatalf(format string, v ...interface{}) {
	l.outputf("FATAL", format, v...)
	os.Exit(1)
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

func (l *paranoidLogger) output(mtype string, v ...interface{}) {
	fmt := "[" + mtype + "] "
	// Add an extra space if the message type (mtype) is only 4 letters long
	if len(mtype) == 4 {
		fmt += " " + l.component + ":"
	} else {
		fmt += l.component + ":"
	}

	var args []interface{}
	args = append(args, fmt)
	args = append(args, v...)

	log.Println(args...)
}

func (l *paranoidLogger) outputf(mtype string, format string, v ...interface{}) {
	fmt := "[" + mtype + "] "
	// Add an extra space if the message type (mtype) is only 4 letters long
	if len(mtype) == 4 {
		fmt += " " + l.component + ": " + format
	} else {
		fmt += l.component + ": " + format
	}

	log.Printf(fmt, v...)
}

func createFileWriter(logPath string, packageName string) (io.Writer, error) {
	return os.OpenFile(path.Join(logPath, packageName+".log"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
}
