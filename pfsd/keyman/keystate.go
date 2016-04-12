package keyman

import (
	"encoding/gob"
	"fmt"
	pb "github.com/cpssd/paranoid/proto/raft"
	"io"
	"os"
	"path"
	"sync"
)

const KSM_FILE_NAME string = "key_state"

var StateMachine *KeyStateMachine

type keyStateElement struct {
	Generation int
	Owner      *pb.Node
	Holder     *pb.Node
}

type KeyStateMachine struct {
	CurrentGeneration    int
	InProgressGeneration int
	// Key is generation number.
	// Value is a list of Node UUID's.
	Nodes map[int]([]string)

	// The first index indicates the generation.
	// The second index is unimportant as order doesn't matter there.
	Elements map[int]([]*keyStateElement)

	stateLock sync.Mutex
	fileLock  sync.Mutex

	// This is, once again, to avoid an import cycle
	PfsDir string
}

func NewKSM(pfsDir string) *KeyStateMachine {
	return &KeyStateMachine{
		CurrentGeneration:    -1,
		InProgressGeneration: -1,
		Nodes:                make(map[int]([]string)),
		Elements:             make(map[int]([]*keyStateElement)),
		PfsDir:               pfsDir,
	}
}

func NewKSMFromReader(reader io.Reader) (*KeyStateMachine, error) {
	ksm := new(KeyStateMachine)
	dec := gob.NewDecoder(reader)
	err := dec.Decode(ksm)
	if err != nil {
		Log.Error("Failed decoding GOB KeyStateMachine data:", err)
		return nil, fmt.Errorf("failed decoding from GOB: %s", err)
	}
	return ksm, nil
}

func NewKSMFromPFSDir(pfsDir string) (*KeyStateMachine, error) {
	file, err := os.Open(path.Join(pfsDir, "meta", KSM_FILE_NAME))
	if err != nil {
		Log.Errorf("Unable to open %s for reading state: %s", pfsDir, err)
		return nil, fmt.Errorf("unable to open %s: %s", pfsDir, err)
	}
	defer file.Close()
	return NewKSMFromReader(file)
}

func (ksm *KeyStateMachine) NewGeneration(generationNumber int, nodeIds []string) error {
	ksm.stateLock.Lock()
	defer ksm.stateLock.Unlock()

	if generationNumber <= ksm.InProgressGeneration {
		return fmt.Errorf("proposed generation too low; given %d, minimum %d", generationNumber, ksm.CurrentGeneration+1)
	}
	if generationNumber > ksm.CurrentGeneration+1 {
		if len(ksm.Elements[ksm.CurrentGeneration+1]) > 0 {
			return fmt.Errorf("generation %d already in progress", ksm.CurrentGeneration+1)
		}
		return fmt.Errorf("generation number too large; next in sequence: %d", ksm.CurrentGeneration+1)
	}
	ksm.Nodes[generationNumber] = nodeIds
	ksm.Elements[generationNumber] = []*keyStateElement{}
	return nil
}

func (ksm KeyStateMachine) NodeInGeneration(generationNumber int, nodeId string) bool {
	nodeIds, ok := ksm.Nodes[generationNumber]
	if !ok {
		return false
	}
	for _, v := range nodeIds {
		if v == nodeId {
			return true
		}
	}
	return false
}

func (ksm *KeyStateMachine) Update(req *pb.KeyStateMessage) error {
	ksm.stateLock.Lock()
	defer ksm.stateLock.Unlock()

	if _, ok := ksm.Elements[int(req.CurrentGeneration)]; !ok {
		return fmt.Errorf("generation %d has not yet been initialised", req.CurrentGeneration)
	}

	elem := &keyStateElement{
		Generation: int(req.CurrentGeneration),
		Owner:      req.GetKeyOwner(),
		Holder:     req.GetKeyHolder(),
	}

	// If a new generation is created, the state machine will only
	// update its CurrentGeneration when enough generation N+1 elements
	// exist for every node in the cluster to unlock if locked.
	var backupGeneration int
	var backupElements []*keyStateElement
	if elem.Generation > ksm.CurrentGeneration && ksm.canUpdateGeneration(elem.Generation) {
		backupGeneration = ksm.CurrentGeneration
		backupElements = ksm.Elements[ksm.CurrentGeneration]
		ksm.CurrentGeneration = elem.Generation
		delete(ksm.Elements, ksm.CurrentGeneration)
		delete(ksm.Nodes, ksm.CurrentGeneration)
	}
	ksm.Elements[elem.Generation] = append(ksm.Elements[elem.Generation], elem)
	err := ksm.SerialiseToPFSDir()
	if err != nil {
		// If the serialisation fails, undo the update.
		ksm.Elements[elem.Generation] = ksm.Elements[elem.Generation][:len(ksm.Elements[elem.Generation])-1]
		ksm.CurrentGeneration = backupGeneration
		ksm.Elements[ksm.CurrentGeneration] = backupElements
		return fmt.Errorf("failed to commit change to KeyStateMachine: %s", err)
	}

	Log.Verbosef("KeyPiece exchange tracked: %s -> %s", elem.Owner.NodeId, elem.Holder.NodeId)
	return nil
}

// Count all of the keys grouped by owner and make sure they meet a minimum.
func (ksm KeyStateMachine) canUpdateGeneration(generation int) bool {
	// Map of UUIDs (as string) to int
	owners := make(map[string]int)
	for _, v := range ksm.Elements[generation] {
		owners[v.Owner.NodeId] += 1
	}
	minNodesRequired := len(ksm.Nodes[generation])/2 + 1
	for _, count := range owners {
		if count < minNodesRequired {
			return false
		}
	}
	return true
}

func (ksm *KeyStateMachine) Serialise(writer io.Writer) error {
	enc := gob.NewEncoder(writer)
	err := enc.Encode(ksm)
	if err != nil {
		Log.Error("Failed encoding KeyStateMachine to GOB:", err)
		return fmt.Errorf("failed encoding KeyStateMachine to GOB:", err)
	}
	return nil
}

func (ksm *KeyStateMachine) SerialiseToPFSDir() error {
	ksm.fileLock.Lock()
	defer ksm.fileLock.Unlock()
	file, err := os.Create(path.Join(ksm.PfsDir, "meta", KSM_FILE_NAME))
	if err != nil {
		Log.Errorf("Unable to open %s for writing state: %s", ksm.PfsDir, err)
		return fmt.Errorf("unable to open %s for writing state: %s", ksm.PfsDir, err)
	}
	defer file.Close()
	return ksm.Serialise(file)
}
