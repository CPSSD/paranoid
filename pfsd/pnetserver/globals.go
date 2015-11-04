// Package pnetserver implements the ParanoidNetwork gRPC server.
// globals.go contains data used by each gRPC handler in pnetserver.
package pnetserver

type ParanoidServer struct{}

// Path to the PFS root directory
var PFSDir string
