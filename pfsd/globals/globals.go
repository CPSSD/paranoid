package globals

import (
	"fmt"
	"github.com/cpssd/paranoid/pfsd/keyman"
	"sync"
	"time"
)

// Node struct
type Node struct {
	IP         string
	Port       string
	CommonName string
	UUID       string
}

func (n Node) String() string {
	return fmt.Sprintf("%s:%s", n.IP, n.Port)
}

// If UPnP is enabled and a port mapping has been establised.
var UPnPEnabled bool

// ResetInterval containing how often the Renew function has to be called
var ResetInterval time.Duration

// DiscoveryAddr contains the connection sting to the discovery server
var DiscoveryAddr string

// Nodes instance which controls all the information about other pfsd instances
var Nodes = nodes{m: make(map[Node]bool)}

var Server string
var Port int

// Common Name of the cert PFSD is using
var CommonName string

// UUID of the node PFSD is managing
var UUID string

// If true, TLS is being used in all connections to and from PFSD
var TLSEnabled bool

// If true, PFSD will not verify the TLS credentials of anything it connects to
var TLSSkipVerify bool

// Global waitgroup for all goroutines in PFSD
var Wait sync.WaitGroup
var Quit = make(chan bool) // Doesn't matter what the channel holds

// --------------------------------------------
// ---- nodes ---- //
// --------------------------------------------

type nodes struct {
	m    map[Node]bool
	lock sync.Mutex
}

func (ns *nodes) Add(n Node) {
	ns.lock.Lock()
	defer ns.lock.Unlock()
	ns.m[n] = true
}

func (ns *nodes) Remove(n Node) {
	ns.lock.Lock()
	defer ns.lock.Unlock()
	delete(ns.m, n)
}

func (ns *nodes) GetAll() []Node {
	ns.lock.Lock()
	defer ns.lock.Unlock()

	var res []Node
	for node := range ns.m {
		res = append(res, node)
	}
	return res
}

//	--------------------
//	---- Encryption ----
//	--------------------

// Global key used by this instance of PFSD
var EncryptionKey keyman.Key

// Indicates when the system has been locked and keys have been distributed
var SystemLocked bool = false

// Map of Nodes to their KeyPiece held by this node
var HeldKeyPieces = make(map[Node]keyman.KeyPiece)
