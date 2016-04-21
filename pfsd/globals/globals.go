package globals

import (
	"crypto/rand"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/logger"
	"github.com/cpssd/paranoid/pfsd/keyman"
	"github.com/cpssd/paranoid/raft"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

const (
	PASSWORD_SALT_LENGTH int = 64
)

var Log *logger.ParanoidLogger

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

type FileSystemAttributes struct {
	Encrypted     bool       `json:"encrypted"`
	KeyGenerated  bool       `json:"keygenerated"`
	NetworkOff    bool       `json:"networkoff"`
	EncryptionKey keyman.Key `json:"encryptionkey,omitempty"` //The encryption key is only saved to file in this manner if networking is turned off
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
var Nodes = nodes{m: make(map[string]Node)}

var NetworkOff bool

// If true, TLS is being used in all connections to and from PFSD
var TLSEnabled bool

// If true, PFSD will not verify the TLS credentials of anything it connects to
var TLSSkipVerify bool

// The hash of the password required to join the raft cluster
var PoolPasswordHash []byte
var PoolPasswordSalt []byte

func SetPoolPasswordHash(password string) error {
	PoolPasswordHash = make([]byte, 0)
	PoolPasswordSalt = make([]byte, PASSWORD_SALT_LENGTH)
	n, err := io.ReadFull(rand.Reader, PoolPasswordSalt)
	if err != nil {
		return err
	}
	if n != PASSWORD_SALT_LENGTH {
		return errors.New("unable to read correct number of bytes from random number generator")
	}

	if password != "" {
		PoolPasswordHash, err = bcrypt.GenerateFromPassword(append(PoolPasswordSalt, []byte(password)...), bcrypt.DefaultCost)
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
	m    map[string]Node
	lock sync.Mutex
}

func (ns *nodes) Add(n Node) {
	ns.lock.Lock()
	defer ns.lock.Unlock()
	ns.m[n.UUID] = n
}

func (ns *nodes) GetNode(uuid string) (Node, error) {
	ns.lock.Lock()
	defer ns.lock.Unlock()

	Log.Info("Getting", uuid)
	for i, _ := range ns.m {
		Log.Info("In map:", i)
	}

	node, ok := ns.m[uuid]
	if !ok {
		return node, errors.New("unrecognised uuid")
	}
	return node, nil
}

func (ns *nodes) Remove(n Node) {
	ns.lock.Lock()
	defer ns.lock.Unlock()
	delete(ns.m, n.UUID)
}

func (ns *nodes) GetAll() []Node {
	ns.lock.Lock()
	defer ns.lock.Unlock()

	var res []Node
	for _, node := range ns.m {
		res = append(res, node)
	}
	return res
}

//	--------------------
//	---- Encryption ----
//	--------------------

// Global key used by this instance of PFSD
var Encrypted bool
var KeyGenerated bool
var EncryptionKey *keyman.Key

// Indicates when the system has been locked and keys have been distributed
var SystemLocked bool = false

var keyPieceStoreLock sync.Mutex

type KeyPieceMap map[string]*keyman.KeyPiece
type KeyPieceStore map[int64]KeyPieceMap

// Returns nil if the piece is not found
func (ks KeyPieceStore) GetPiece(generation int64, nodeUUID string) *keyman.KeyPiece {
	keyPieceStoreLock.Lock()
	defer keyPieceStoreLock.Unlock()
	keymap, ok := ks[generation]
	if !ok {
		return nil
	}

	piece, ok := keymap[nodeUUID]
	if !ok {
		return nil
	}
	return piece
}

func (ks KeyPieceStore) AddPiece(generation int64, nodeUUID string, piece *keyman.KeyPiece) error {
	keyPieceStoreLock.Lock()
	defer keyPieceStoreLock.Unlock()

	_, ok := ks[generation]
	if !ok {
		ks[generation] = make(KeyPieceMap)
	}

	ks[generation][nodeUUID] = piece
	return ks.SaveToDisk()
}

func (ks KeyPieceStore) DeletePiece(generation int64, nodeUUID string) error {
	keyPieceStoreLock.Lock()
	defer keyPieceStoreLock.Unlock()

	_, ok := ks[generation]
	if !ok {
		return nil
	}

	delete(ks[generation], nodeUUID)
	return ks.SaveToDisk()
}

func (ks KeyPieceStore) DeleteGeneration(generation int64) error {
	keyPieceStoreLock.Lock()
	defer keyPieceStoreLock.Unlock()
	delete(ks, generation)
	return ks.SaveToDisk()
}

func (ks KeyPieceStore) SaveToDisk() error {
	piecePath := path.Join(ParanoidDir, "meta", "pieces-new")
	file, err := os.Create(piecePath)
	if err != nil {
		Log.Errorf("Unable to open %s for storing pieces: %s", piecePath, file)
		return fmt.Errorf("Unable to open %s for storing pieces: %s", piecePath, file)
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	err = enc.Encode(ks)
	if err != nil {
		Log.Error("Failed encoding KeyPieceStore to GOB:", err)
		return fmt.Errorf("failed encoding KeyPieceStore to GOB: %s", err)
	}
	err = os.Rename(piecePath, path.Join(ParanoidDir, "meta", "pieces"))
	if err != nil {
		Log.Error("Failed to save KeyPieceStore to file:", err)
		return fmt.Errorf("Failed to save KeyPieceStore to file: %s", err)
	}
	return nil
}

var HeldKeyPieces = make(KeyPieceStore)
