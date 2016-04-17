package pnetclient

import (
	"fmt"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/keyman"
)

// Chunks key and sends the pieces to other nodes on the network.
func Distribute(key *keyman.Key, peers []*globals.Node, generation int) error {
	numPieces := int64(len(peers) + 1)
	requiredPieces := numPieces/2 + 1
	Log.Info("Generating pieces.")
	pieces, err := keyman.GeneratePieces(globals.EncryptionKey, numPieces, requiredPieces)
	if err != nil {
		Log.Error("Could not chunk key:", err)
		return fmt.Errorf("could not chunk key: %s", err)
	}
	// We always keep the first piece and distribute the rest
	globals.HeldKeyPieces.AddPiece(globals.ThisNode.UUID, pieces[0])

	for i := 1; i < len(pieces); i++ {
		SendKeyPiece(pieces[i])
	}
	return nil
}
