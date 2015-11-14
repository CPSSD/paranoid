// Package globals contains the shared variables and functions in pfsd
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

// Nodes class containing all methods necessary to maniplate the nodes
var Nodes = nodes{m: make(map[Node]bool)}

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

func (ns *nodes) Delete(n Node) {
	ns.Remove(n)
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
