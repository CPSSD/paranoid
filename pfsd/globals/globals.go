package globals

import (
	"sync"
	"time"
)

// Node struct
type Node struct {
	IP   string
	Port string
}

// ResetInterval containing how often the Renew function has to be called
var ResetInterval time.Duration

// DiscoveryAddr contains the connection sting to the discovery server
var DiscoveryAddr string

// Nodes instance which controls all the information about other pfsd instances
var Nodes = nodes{m: make(map[Node]bool)}

var Port int

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
