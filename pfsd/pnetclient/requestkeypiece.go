package pnetclient

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/keyman"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"math/big"
)

// A struct which is either a pointer to a KeyPiece or an error.
// Basically like a union from C.
type keyPieceUnion struct {
	piece *keyman.KeyPiece
	err   error
}

func RequestKeyPiece(uuid string) (*keyman.KeyPiece, error) {
	node, err := globals.Nodes.GetNode(uuid)
	if err != nil {
		return nil, errors.New("could not find node details")
	}

	conn, err := Dial(node)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %s", node, err)
	}
	defer conn.Close()

	client := pb.NewParanoidNetworkClient(conn)

	thisNodeProto := &pb.Node{
		Ip:         globals.ThisNode.IP,
		Port:       globals.ThisNode.Port,
		CommonName: globals.ThisNode.CommonName,
		Uuid:       globals.ThisNode.UUID,
	}
	pieceProto, err := client.RequestKeyPiece(context.Background(), thisNodeProto)
	if err != nil {
		Log.Warn("Failed requesting KeyPiece from", node, "Error:", err)
		return nil, fmt.Errorf("failed requesting KeyPiece from %s: %s", node, err)
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
	return piece, nil
}
