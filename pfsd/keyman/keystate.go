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

	lock sync.Mutex

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
	ksm.lock.Lock()
	defer ksm.lock.Unlock()

	if generationNumber <= ksm.InProgressGeneration {
		return fmt.Errorf("proposed generation too low; given %d, minimum %d", generationNumber, ksm.CurrentGeneration+1)
	}
	if generationNumber > ksm.CurrentGeneration+1 && len(ksm.Elements[ksm.CurrentGeneration+1]) > 0 {
		return fmt.Errorf("generation %d already in progress", ksm.CurrentGeneration+1)
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
	ksm.lock.Lock()
	defer ksm.lock.Unlock()

	if _, ok := ksm.Elements[int(req.CurrentGeneration)]; !ok {
		return fmt.Errorf("generation %d has not yet been initialised", req.CurrentGeneration)
	}

	elem := &keyStateElement{
		Generation: int(req.CurrentGeneration),
		Owner:      req.GetKeyOwner(),
		Holder:     req.GetKeyHolder(),
	}

	ksm.Elements[elem.Generation] = append(ksm.Elements[elem.Generation], elem)
	// If a new generation is created, the state machine will only
	// update its CurrentGeneration when enough generation N+1 elements
	// exist for every node in the cluster to unlock if locked.
	var backupGeneration int
	var backupElements []*keyStateElement
	var backupNodes []string
	if elem.Generation > ksm.CurrentGeneration && ksm.canUpdateGeneration(elem.Generation) {
		backupGeneration = ksm.CurrentGeneration
		backupElements = ksm.Elements[ksm.CurrentGeneration]
		backupNodes = ksm.Nodes[ksm.CurrentGeneration]
		delete(ksm.Elements, ksm.CurrentGeneration)
		delete(ksm.Nodes, ksm.CurrentGeneration)
		ksm.CurrentGeneration = elem.Generation
	}
	err := ksm.SerialiseToPFSDir()
	if err != nil {
		// If the serialisation fails, undo the update.
		ksm.Elements[elem.Generation] = ksm.Elements[elem.Generation][:len(ksm.Elements[elem.Generation])-1]
		ksm.CurrentGeneration = backupGeneration
		ksm.Elements[ksm.CurrentGeneration] = backupElements
		ksm.Nodes[ksm.CurrentGeneration] = backupNodes
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
	if len(owners) != len(ksm.Elements[generation]) {
		return false
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
	ksmpath := path.Join(ksm.PfsDir, "meta", KSM_FILE_NAME)
	file, err := os.Create(ksmpath + "-new")
	if err != nil {
		Log.Errorf("Unable to open %s for writing state: %s", ksm.PfsDir, err)
		return fmt.Errorf("unable to open %s for writing state: %s", ksm.PfsDir, err)
	}
	err = ksm.Serialise(file)
	file.Close()
	if err == nil {
		err := os.Rename(ksmpath+"-new", ksmpath)
		if err != nil {
			// We ignore the following error because if both file operations fail they are very
			// likely caused by the same thing, so one error will give information for both.
			os.Remove(ksmpath + "-new")
			return fmt.Errorf("unable to rename %s to %s: %s", ksmpath+"-new", ksmpath, err)
		}
	}
	return err
}
