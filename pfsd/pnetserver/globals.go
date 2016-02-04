// Package pnetserver implements the ParanoidNetwork gRPC server.
// globals.go contains data used by each gRPC handler in pnetserver.
package pnetserver

import (
	"github.com/cpssd/paranoid/logger"
)

type ParanoidServer struct{}

// Path to the PFS root directory
var ParanoidDir string

var Log *logger.ParanoidLogger
