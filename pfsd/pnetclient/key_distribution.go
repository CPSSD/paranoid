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
		SendKeyPiece(peers[i-1].UUID, int64(generation), pieces[i])
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
		replicationLoop:
			for {
				select {
				case <-ksm.Events:
					//Keep this channel clear
				default:
					done := true
					for g := ksm.GetCurrentGeneration(); g <= ksm.GetInProgressGenertion(); g++ {
						if !ksm.NeedsReplication(globals.ThisNode.UUID, g) {
							continue
						}
						done = false
						nodes, err := ksm.GetNodes(g)
						if err != nil {
							Log.Warn("Unable to get nodes for generation", g, ":", err)
							continue
						}
						var peers []globals.Node
						for _, v := range nodes {
							if v != globals.ThisNode.UUID {
								globalNode, err := globals.Nodes.GetNode(v)
								if err != nil {
									Log.Errorf("Unable to lookup node %s: %s", v, err)
								} else {
									peers = append(peers, globalNode)
								}
							}
						}
						Distribute(globals.EncryptionKey, peers, int(g))
					}
					for i := int64(0); i < ksm.GetCurrentGeneration(); i++ {
						globals.HeldKeyPieces.DeleteGeneration(i)
					}
					if done {
						break replicationLoop
					}
				}
			}
		}
	}
}
