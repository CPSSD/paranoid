package pnetserver

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/keyman"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	raftpb "github.com/cpssd/paranoid/proto/raft"
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
			raftOwner := &raftpb.Node{
				Ip:         req.OwnerNode.Ip,
				Port:       req.OwnerNode.Port,
				CommonName: req.OwnerNode.CommonName,
				NodeId:     req.OwnerNode.Uuid,
			}
			raftHolder := &raftpb.Node{
				Ip:         globals.ThisNode.IP,
				Port:       globals.ThisNode.Port,
				CommonName: globals.ThisNode.CommonName,
				NodeId:     globals.ThisNode.UUID,
			}
			_, err := globals.RaftNetworkServer.RequestKeyStateUpdate(raftOwner, raftHolder,
				int64(keyman.StateMachine.CurrentGeneration+1))
			if err != nil {
				return &pb.EmptyMessage{}, grpc.Errorf(codes.FailedPrecondition, "failed to commit to Raft: %s", err)
			}
			globals.HeldKeyPieces.AddPiece(node, piece)
			Log.Info("Received KeyPiece from", node)
			return &pb.EmptyMessage{}, nil
		}
	}
	Log.Warn("OwnerNode not found:", req.OwnerNode)
	err := grpc.Errorf(codes.NotFound, "OwnerNode not found in local database: %v", req.OwnerNode)
	return &pb.EmptyMessage{}, err
}
