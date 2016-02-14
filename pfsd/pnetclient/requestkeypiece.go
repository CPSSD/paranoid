package pnetclient

import (
	"fmt"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/keyman"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"math/big"
)

func RequestKeyPieces() ([]*keyman.KeyPiece, error) {
	var pieces []*keyman.KeyPiece
	for _, node := range globals.Nodes.GetAll() {
		conn, err := Dial(node)
		if err != nil {
			Log.Error("RequestKeyPiece: failed to dial ", node)
			return nil, fmt.Errorf("failed to dial %s", node)
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
			Log.Error("Failed requesting KeyPiece from", node, "Error:", err)
		}
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
		pieces = append(pieces, piece)
	}
	return pieces, nil
}
