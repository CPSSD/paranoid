package pnetserver

import (
	"github.com/cpssd/paranoid/libpfs/commands"
	"github.com/cpssd/paranoid/libpfs/returncodes"
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
		if node.IP == req.OwnerNode.Ip && node.Port == req.OwnerNode.Port && node.CommonName == req.OwnerNode.CommonName {
			prime := *big.Int
			prime.SetBytes(req.Prime)
			piece := keyman.KeyPiece{
				Data:              req.Data,
				ParentFingerprint: req.ParentFingerprint,
				Prime:             prime,
				Seq:               req.Seq,
			}
			globals.HeldKeyPieces[node] = piece
			return &pb.EmptyMessage{}, nil
		}
	}
	err := grpc.Errorf(codes.NotFound, "OwnerNode not found in local database: %v", req.OwnerNode)
	return &pb.EmptyMessage{}, err
}
