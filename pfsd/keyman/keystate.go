package keyman

import (
	"encoding/gob"
	"fmt"
	pb "github.com/cpssd/paranoid/proto/raft"
	"io"
)

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
	Elements [][]*keyStateElement
}

func NewKSM() *KeyStateMachine {
	return &KeyStateMachine{
		CurrentGeneration: -1,
		Elements:          make([][]*keyStateElement, 0, 256),
	}
}

func NewKSMFromReader(reader io.Reader) (*KeyStateMachine, error) {
	ksm := new(KeyStateMachine)
	dec := gob.NewDecoder(reader)
	err := dec.Decode(ksm)
	if err != nil {
		Log.Fatal("Failed decoding GOB KeyStateMachine data:", err)
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
	if elem.generation > ksm.CurrentGeneration {
		ksm.CurrentGeneration = elem.generation
		if elem.generation > cap(ksm.Elements) {
			tmp := make([][]*keyStateElement, len(ksm.Elements), cap(ksm.Elements)*2)
			ksm.Elements = tmp
		}
	}
	ksm.Elements[elem.generation] = append(ksm.Elements[elem.generation], elem)

	return nil
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
