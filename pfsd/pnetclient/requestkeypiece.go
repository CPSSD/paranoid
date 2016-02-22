package pnetclient

import (
	"fmt"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/keyman"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"math/big"
	"sync"
)

// A struct which is either a pointer to a KeyPiece or an error.
// Basically like a union from C.
type keyPieceUnion struct {
	piece *keyman.KeyPiece
	err   error
}

func requestKeyPiece(node globals.Node, c chan keyPieceUnion, w *sync.WaitGroup) {
	defer w.Done()
	conn, err := Dial(node)
	if err != nil {
		c <- keyPieceUnion{
			piece: nil,
			err:   fmt.Errorf("failed to dial %s: %s", node, err),
		}
	}
	defer conn.Close()

	client := pb.NewParanoidNetworkClient(conn)

	thisNodeProto := &pb.PingRequest{
		Ip:         globals.ThisNode.IP,
		Port:       globals.ThisNode.Port,
		CommonName: globals.ThisNode.CommonName,
	}
	pieceProto, err := client.RequestKeyPiece(context.Background(), thisNodeProto)
	if err != nil {
		Log.Warn("Failed requesting KeyPiece from", node, "Error:", err)
		c <- keyPieceUnion{
			piece: nil,
			err:   fmt.Errorf("failed requesting KeyPiece from %s: %s", node, err),
		}
		return
	}
	Log.Info("Received KeyPiece from", node)
	var fingerprintArray [32]byte
	copy(fingerprintArray[:], pieceProto.ParentFingerprint)
	var primeBig big.Int
	primeBig.SetBytes(pieceProto.Prime)
	piece := &keyman.KeyPiece{
		Data:              pieceProto.Data,
		ParentFingerprint: fingerprintArray,
		Prime:             &primeBig,
		Seq:               pieceProto.Seq,
	}
	c <- keyPieceUnion{
		piece: piece,
		err:   nil,
	}
	return
}

// Asks every known node for a KeyPiece belonging to this node.
func RequestKeyPieces() ([]*keyman.KeyPiece, error) {
	var pieces []*keyman.KeyPiece
	unionChan := make(chan keyPieceUnion, len(globals.Nodes.GetAll()))
	var wait sync.WaitGroup

	for _, node := range globals.Nodes.GetAll() {
		wait.Add(1)
		go requestKeyPiece(node, unionChan, &wait)
	}
	Log.Info("Waiting for KeyPieces.")
	wait.Wait()
	close(unionChan)

	for union := range unionChan {
		if union.err != nil {
			Log.Warn(union.err)
		} else {
			pieces = append(pieces, union.piece)
		}
	}
	return pieces, nil
}
