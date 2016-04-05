package globals

import (
	"encoding/gob"
	"fmt"
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/pfsd/keyman"
	"github.com/cpssd/paranoid/raft"
	"golang.org/x/crypto/bcrypt"
	"os"
	"path"
	"sync"
	"time"
)

var Log *logger.ParanoidLogger

// Node struct
type Node struct {
	IP         string
	Port       string
	CommonName string
	UUID       string
}

type FileSystemAttributes struct {
	Encrypted     bool       `json:"encrypted"`
	KeyGenerated  bool       `json:"keygenerated"`
	NetworkOff    bool       `json:"networkoff"`
	EncryptionKey keyman.Key `json:"encryptionkey,omitempty"` //The encryption key is only saved to file in this manner if networking is turned off
}

func (n Node) String() string {
	return fmt.Sprintf("%s:%s", n.IP, n.Port)
}

var RaftNetworkServer *raft.RaftNetworkServer

var ParanoidDir string
var MountPoint string

// Time at which PFSD started. Used for calculating uptime.
var BootTime time.Time

// Node information for the current node
var ThisNode Node

// If UPnP is enabled and a port mapping has been establised.
var UPnPEnabled bool

// ResetInterval containing how often the Renew function has to be called
var ResetInterval time.Duration

// DiscoveryAddr contains the connection sting to the discovery server
var DiscoveryAddr string

// Nodes instance which controls all the information about other pfsd instances
var Nodes = nodes{m: make(map[Node]bool)}

// If true, TLS is being used in all connections to and from PFSD
var TLSEnabled bool

// If true, PFSD will not verify the TLS credentials of anything it connects to
var TLSSkipVerify bool

// The hash of the password required to join the raft cluster
var PoolPasswordHash []byte

func SetPoolPasswordHash(password string) error {
	PoolPasswordHash = make([]byte, 0)
	if password != "" {
		var err error
		PoolPasswordHash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return err
	}
	return nil
}

// Global waitgroup for all goroutines in PFSD
var Wait sync.WaitGroup
var Quit = make(chan bool) // Doesn't matter what the channel holds
var ShuttingDown bool

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
var Encrypted bool
var EncryptionKey *keyman.Key

// Indicates when the system has been locked and keys have been distributed
var SystemLocked bool = false

var keyPieceStoreLock sync.Mutex

type KeyPieceStore map[Node]*keyman.KeyPiece

// Returns nil if the piece is not found
func (ks KeyPieceStore) GetPiece(node Node) *keyman.KeyPiece {
	keyPieceStoreLock.Lock()
	defer keyPieceStoreLock.Unlock()
	piece, ok := ks[node]
	if !ok {
		return nil
	}
	return piece
}

func (ks KeyPieceStore) AddPiece(node Node, piece *keyman.KeyPiece) {
	keyPieceStoreLock.Lock()
	defer keyPieceStoreLock.Unlock()
	ks[node] = piece
	ks.SaveToDisk()
}

func (ks KeyPieceStore) DeletePiece(node Node) {
	keyPieceStoreLock.Lock()
	defer keyPieceStoreLock.Unlock()
	delete(ks, node)
	ks.SaveToDisk()
}

func (ks KeyPieceStore) SaveToDisk() {
	piecePath := path.Join(ParanoidDir, "meta", "pieces")
	file, err := os.Create(piecePath)
	if err != nil {
		Log.Errorf("Unable to open %s for storing pieces: %s", piecePath, file)
		return
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	err = enc.Encode(ks)
	if err != nil {
		Log.Error("Failed encoding KeyPieceStore to GOB:", err)
	}
}

var HeldKeyPieces = make(KeyPieceStore)
