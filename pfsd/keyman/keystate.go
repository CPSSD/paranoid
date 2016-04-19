package keyman

import (
	"encoding/gob"
	"errors"
	"fmt"
	pb "github.com/cpssd/paranoid/proto/raft"
	"io"
	"os"
	"path"
	"sync"
)

const KSM_FILE_NAME string = "key_state"

var ErrGenerationDeprecated = errors.New("given generation was created before the current generation was set")
var StateMachine *KeyStateMachine

type keyStateElement struct {
	Owner  *pb.Node
	Holder *pb.Node
}

type Generation struct {
	//A list of all nodes included in the generation
	Nodes    []*pb.Node
	Elements []*keyStateElement
}

func (g *Generation) AddElement(elem *keyStateElement) {
	g.Elements = append(g.Elements, elem)
}

func (g *Generation) RemoveElement() {
	g.Elements = g.Elements[:len(g.Elements)-1]
}

type KeyStateMachine struct {
	CurrentGeneration    int64
	InProgressGeneration int64
	DeprecatedGeneration int64

	// The first index indicates the generation.
	// The second index is unimportant as order doesn't matter there.
	Generations map[int64]*Generation

	lock sync.Mutex

	// This is, once again, to avoid an import cycle
	PfsDir string
}

func NewKSM(pfsDir string) *KeyStateMachine {
	return &KeyStateMachine{
		CurrentGeneration:    -1,
		InProgressGeneration: -1,
		DeprecatedGeneration: -1,
		Generations:          make(map[int64]*Generation),
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

func (ksm *KeyStateMachine) UpdateFromStateFile(filePath string) error {
	ksm.lock.Lock()
	defer ksm.lock.Unlock()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open %s: %s", filePath, err)
	}
	defer file.Close()
	tmpKSM, err := NewKSMFromReader(file)
	if err != nil {
		return fmt.Errorf("unable to create new key state machine: %s", err)
	}

	ksm.CurrentGeneration = tmpKSM.CurrentGeneration
	ksm.InProgressGeneration = tmpKSM.InProgressGeneration
	ksm.DeprecatedGeneration = tmpKSM.DeprecatedGeneration
	ksm.Generations = tmpKSM.Generations
	return nil
}

func (ksm *KeyStateMachine) NewGeneration(newNode *pb.Node) (generationNumber int64, err error) {
	ksm.lock.Lock()
	defer ksm.lock.Unlock()

	var existingNodes []*pb.Node
	if gen, ok := ksm.Generations[ksm.CurrentGeneration]; ok {
		existingNodes = gen.Nodes
	}
	ksm.InProgressGeneration++
	ksm.Generations[ksm.InProgressGeneration] = &Generation{
		Nodes:    append(existingNodes, newNode),
		Elements: []*keyStateElement{},
	}

	err = ksm.SerialiseToPFSDir()
	if err != nil {
		delete(ksm.Generations, ksm.InProgressGeneration)
		ksm.InProgressGeneration--
		return 0, err
	}
	return ksm.InProgressGeneration, nil
}

func (ksm KeyStateMachine) NodeInGeneration(generationNumber int64, nodeId string) bool {
	generation, ok := ksm.Generations[generationNumber]
	if !ok {
		return false
	}
	for _, v := range generation.Nodes {
		if v.NodeId == nodeId {
			return true
		}
	}
	return false
}

func (ksm *KeyStateMachine) Update(req *pb.KeyStateMessage) error {
	ksm.lock.Lock()
	defer ksm.lock.Unlock()

	if req.Generation != ksm.CurrentGeneration && req.Generation <= ksm.DeprecatedGeneration {
		return ErrGenerationDeprecated
	}

	if _, ok := ksm.Generations[req.Generation]; !ok {
		return fmt.Errorf("generation %d has not yet been initialised", req.Generation)
	}

	elem := &keyStateElement{
		Owner:  req.GetKeyOwner(),
		Holder: req.GetKeyHolder(),
	}

	ksm.Generations[req.Generation].AddElement(elem)
	// If a new generation is created, the state machine will only
	// update its CurrentGeneration when enough generation N+1 elements
	// exist for every node in the cluster to unlock if locked.
	var backupGeneration int64
	var backupDeprecatedGen int64
	var backupGenerations map[int64]*Generation
	updatedGeneration := false
	if req.Generation > ksm.CurrentGeneration && ksm.canUpdateGeneration(req.Generation) {
		updatedGeneration = true
		backupGeneration = ksm.CurrentGeneration
		backupDeprecatedGen = ksm.DeprecatedGeneration
		backupGenerations = ksm.Generations

		ksm.Generations = make(map[int64]*Generation)
		ksm.InProgressGeneration++
		ksm.DeprecatedGeneration = ksm.InProgressGeneration
		ksm.CurrentGeneration = req.Generation
		ksm.Generations[req.Generation] = backupGenerations[req.Generation]
	}
	err := ksm.SerialiseToPFSDir()
	if err != nil {
		// If the serialisation fails, undo the update.
		if updatedGeneration {
			ksm.CurrentGeneration = backupGeneration
			ksm.InProgressGeneration--
			ksm.DeprecatedGeneration = backupDeprecatedGen
			ksm.Generations = backupGenerations
		}
		ksm.Generations[req.Generation].RemoveElement()
		return fmt.Errorf("failed to commit change to KeyStateMachine: %s", err)
	}

	Log.Verbosef("KeyPiece exchange tracked: %s -> %s", elem.Owner.NodeId, elem.Holder.NodeId)
	return nil
}

// Count all of the keys grouped by owner and make sure they meet a minimum.
func (ksm KeyStateMachine) canUpdateGeneration(generation int64) bool {
	// Map of UUIDs (as string) to int
	owners := make(map[string]int)
	for _, v := range ksm.Generations[generation].Elements {
		owners[v.Owner.NodeId] += 1
	}
	if len(owners) != len(ksm.Generations[generation].Nodes) {
		return false
	}
	minNodesRequired := len(ksm.Generations[generation].Nodes)/2 + 1
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
