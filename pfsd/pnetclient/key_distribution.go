package pnetclient

import (
	"fmt"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/keyman"
)

// Chunks key and sends the pieces to other nodes on the network.
func Distribute(key *keyman.Key, peers []globals.Node, generation int) error {
	numPieces := int64(len(peers) + 1)
	requiredPieces := numPieces/2 + 1
	Log.Info("Generating pieces.")
	pieces, err := keyman.GeneratePieces(key, numPieces, requiredPieces)
	if err != nil {
		Log.Error("Could not chunk key:", err)
		return fmt.Errorf("could not chunk key: %s", err)
	}
	// We always keep the first piece and distribute the rest
	globals.HeldKeyPieces.AddPiece(int64(generation), globals.ThisNode.UUID, pieces[0])

	for i := 1; i < len(pieces); i++ {
		SendKeyPiece(pieces[i])
	}
	return nil
}

func KSMObserver(ksm *keyman.KeyStateMachine) {
	defer globals.Wait.Done()
	for {
		select {
		case _, ok := <-globals.Quit:
			if !ok {
				return
			}
		case <-ksm.Events:
			for g := ksm.CurrentGeneration; g <= ksm.InProgressGeneration; g++ {
				nodes := make([]globals.Node, len(ksm.Generations[g].Nodes))
				for i, v := range ksm.Generations[g].Nodes {
					globalNode, err := globals.Nodes.GetNode(v)
					if err != nil {
						Log.Errorf("Unable to lookup node %s: %s", v, err)
					}
					nodes[i] = globalNode
				}
				Distribute(globals.EncryptionKey, nodes, int(g))
			}
			for i := int64(0); i < ksm.CurrentGeneration; i++ {
				globals.HeldKeyPieces.DeleteGeneration(i)
			}
		}
	}
}
