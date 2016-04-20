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

func (s *ParanoidServer) SendKeyPiece(ctx context.Context, req *pb.KeyPiece) (*pb.SendKeyPieceResponse, error) {
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

			err := globals.HeldKeyPieces.AddPiece(0, node.UUID, piece)
			if err != nil {
				return &pb.SendKeyPieceResponse{}, grpc.Errorf(codes.FailedPrecondition, "failed to save key piece to disk: %s", err)
			}
			Log.Info("Received KeyPiece from", node)
			if globals.RaftNetworkServer != nil {
				err := globals.RaftNetworkServer.RequestKeyStateUpdate(raftOwner, raftHolder,
					int64(keyman.StateMachine.CurrentGeneration+1))
				if err != nil {
					return &pb.SendKeyPieceResponse{}, grpc.Errorf(codes.FailedPrecondition, "failed to commit to Raft: %s", err)
				}
			} else {
				return &pb.SendKeyPieceResponse{true}, nil
			}
			return &pb.SendKeyPieceResponse{}, nil
		}
	}
	Log.Warn("OwnerNode not found:", req.OwnerNode)
	err := grpc.Errorf(codes.NotFound, "OwnerNode not found in local database: %v", req.OwnerNode)
	return &pb.SendKeyPieceResponse{}, err
}
