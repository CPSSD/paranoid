Paranoid
========

---

This repo is no longer maintained. The project has a new home: https://github.com/pp2p/paranoid

---

## Build Instructions ##

To run the unit tests recursively for the entire project, run `go test ./...` from this directory.

To run the integration tests recursively for the entire project run `go test ./... -tags=integration` from this directory. 

To build a specific binary, consult the README file for that directory.

## Hosted Discovery Server ##
There is a discovery server running at `paranoid.discovery.razoft.net:10101` which can be used to avoid running a discovery server locally while testing. 
