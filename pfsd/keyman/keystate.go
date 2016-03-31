package keyman

import (
	"encoding/gob"
	"fmt"
	pb "github.com/cpssd/paranoid/proto/raft"
	"io"
	"os"
	"path"
)

// Hardcoded because the KSM will not track joined nodes until next sprint (Sprint 7)
const MIN_SHARES_REQUIRED int = 2

var StateMachine *KeyStateMachine

// We have to recreate globals.Node here to avoid an import cycle.
// This just in: Go is officially worse than C at this.
type Node struct {
	IP         string
	Port       string
	CommonName string
	UUID       string
}

func (n Node) String() string {
	return fmt.Sprintf("%s:%s", n.IP, n.Port)
}

type keyStateElement struct {
	generation int
	owner      *Node
	holder     *Node
}

type KeyStateMachine struct {
	CurrentGeneration int

	// The first index indicates the generation.
	// The second index is unimportant as order doesn't matter there.
	Elements map[int]([]*keyStateElement)
}

func NewKSM() *KeyStateMachine {
	return &KeyStateMachine{
		CurrentGeneration: -1,
		Elements:          make(map[int]([]*keyStateElement)),
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
	file, err := os.Open(path.Join(pfsDir, "meta", "key_state"))
	if err != nil {
		Log.Errorf("Unable to open %s for reading state: %s", pfsDir, err)
		return nil, fmt.Errorf("unable to open %s: %s", pfsDir, err)
	}
	defer file.Close()
	ksm := new(KeyStateMachine)
	dec := gob.NewDecoder(file)
	err = dec.Decode(ksm)
	if err != nil {
		Log.Error("Failed decoding GOB KeyStateMachine data:", err)
		return nil, fmt.Errorf("failed decoding from GOB: %s", err)
	}
	return ksm, nil
}

func (ksm *KeyStateMachine) Update(req *pb.KeyStateMessage) error {
	owner := &Node{
		IP:         req.GetKeyOwner().Ip,
		Port:       req.GetKeyOwner().Port,
		CommonName: req.GetKeyOwner().CommonName,
		UUID:       req.GetKeyOwner().NodeId,
	}
	holder := &Node{
		IP:         req.GetKeyHolder().Ip,
		Port:       req.GetKeyHolder().Port,
		CommonName: req.GetKeyHolder().CommonName,
		UUID:       req.GetKeyHolder().NodeId,
	}
	elem := &keyStateElement{
		generation: int(req.CurrentGeneration),
		owner:      owner,
		holder:     holder,
	}

	// If a new generation is created, the state machine will only
	// update its CurrentGeneration when enough generation N+1 elements
	// exist for every node in the cluster to unlock if locked.
	if elem.generation > ksm.CurrentGeneration && ksm.canUpdateGeneration(elem.generation) {
		ksm.CurrentGeneration = elem.generation
		delete(ksm.Elements, elem.generation)
	}
	ksm.Elements[elem.generation] = append(ksm.Elements[elem.generation], elem)

	Log.Verbosef("KeyPiece exchange tracked: %s -> %s", owner.UUID, holder.UUID)
	return nil
}

// Count all of the keys grouped by owner and make sure they meet a minimum.
func (ksm KeyStateMachine) canUpdateGeneration(generation int) bool {
	// Map of UUIDs (as string) to int
	owners := make(map[string]int)
	for _, v := range ksm.Elements[generation] {
		owners[v.owner.UUID] += 1
	}
	for _, count := range owners {
		if count < MIN_SHARES_REQUIRED {
			return false
		}
	}
	return true
}

func (ksm KeyStateMachine) Serialise(writer io.Writer) error {
	enc := gob.NewEncoder(writer)
	err := enc.Encode(ksm)
	if err != nil {
		Log.Error("Failed encoding KeyStateMachine to GOB:", err)
		return fmt.Errorf("failed encoding to GOB: %s", err)
	}
	return nil
}
