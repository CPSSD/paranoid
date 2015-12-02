logger(3) - paranoid standard logging
=====================================

## DESCRIPTION ##
logger (github.com/cpssd/paranoid/logger) allows to log messages in a standard format throughout the project.
It supports writing to stderr, log file, or both. It also allows you to add in your own writers.

To get the documentation for the logger run `godoc -http :6060` and visit  
> http://localhost:6060/pkg/github.com/cpssd/paranoid/logger/

## USAGE ##
```go
// Import the logger
import(
  "github.com/cpssd/paranoid/logger"
)

// Initialize an instance of logger
log := logger.New("currentPackage", "component", "/home/.pfs/example-pfs/meta/log")

// OPTIONAL: Set the output, default: stderr
log.SetOutput(logger.STDERR | logger.LOGFILE) // Prints to both stderr and a log file located at {LOGPATH}/{COMPONENT}.log

// OPTIONAL: Set Logging level, default: INFO
log.SetLogLevel(logger.WARNING)

// OPTIONAL: Add Custom writer
log.AddAdditionalWriter(os.Stdout)
// There is no specific way to remove a writer, so just call log.SetOutput() again

// Logging functions
log.Debug("debug message")      // Only works with LogLevel of DEBUG
log.Verbose("verbose message")  // Only works with LogLevel of VERBOSE or lower importance
log.Info("info message")        // Only works with LogLevel of INFO or lower importance
log.Warn("warning message")     // Only works with LogLevel of WARNING or lower importance
log.Error("error message")      // Only works with LogLevel of ERROR or lower importance
log.Fatal("fatal message")      // Works regardless of the LogLevel set. Quits the program with exit code 1

```

### LogLevel Codes ###
```go
const (
  DEBUG
  VERBOSE
  INFO
  WARNING
  ERROR
)
```

### Output Codes ###
```go
const (
  FILE
  STDERR
)
```

### SetOutput ###
Function SetOutput Accepts output codes OR'd together  
Example: `log.SetOutput(logger.LOGFILE | logger.STDERR)`
