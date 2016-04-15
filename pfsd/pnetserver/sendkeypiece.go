package pnetserver

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/keyman"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"math/big"
)

func (s *ParanoidServer) SendKeyPiece(ctx context.Context, req *pb.KeyPiece) (*pb.EmptyMessage, error) {
	for _, node := range globals.Nodes.GetAll() {
		if node.UUID == req.OwnerNode.Uuid {
			var prime big.Int
			prime.SetBytes(req.Prime)
			// We must convert a slice to an array
			var fingerArray [32]byte
			copy(fingerArray[:], req.ParentFingerprint)
			piece := &keyman.KeyPiece{
				Data:              req.Data,
				ParentFingerprint: fingerArray,
				Prime:             &prime,
				Seq:               req.Seq,
			}
			globals.HeldKeyPieces.AddPiece(node.UUID, piece)
			Log.Info("Received KeyPiece from", node)
			return &pb.EmptyMessage{}, nil
		}
	}
	Log.Warn("OwnerNode not found:", req.OwnerNode)
	err := grpc.Errorf(codes.NotFound, "OwnerNode not found in local database: %v", req.OwnerNode)
	return &pb.EmptyMessage{}, err
}
