package pnetserver

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (s *ParanoidServer) RequestKeyPiece(ctx context.Context, req *pb.PingRequest) (*pb.KeyPiece, error) {
	for _, node := range globals.Nodes.GetAll() {
		if node.UUID == req.Uuid {
			key := globals.HeldKeyPieces.GetPiece(0, node.UUID)
			if key == nil {
				Log.Warn("Key not found for node", node)
				err := grpc.Errorf(codes.NotFound, "Key not found for node %v", node)
				return &pb.KeyPiece{}, err
			}
			keyProto := &pb.KeyPiece{
				Data:              key.Data,
				ParentFingerprint: key.ParentFingerprint[:],
				Prime:             key.Prime.Bytes(),
				Seq:               key.Seq,
			}
			return keyProto, nil
		}
	}

	Log.Warn("Node not found:", req)
	err := grpc.Errorf(codes.NotFound, "Node not found in local database: %v", req)
	return &pb.KeyPiece{}, err
}
