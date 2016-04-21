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

func (s *ParanoidServer) SendKeyPiece(ctx context.Context, req *pb.KeyPieceSend) (*pb.SendKeyPieceResponse, error) {
	var prime big.Int
	prime.SetBytes(req.Key.Prime)
	// We must convert a slice to an array
	var fingerArray [32]byte
	copy(fingerArray[:], req.Key.ParentFingerprint)
	piece := &keyman.KeyPiece{
		Data:              req.Key.Data,
		ParentFingerprint: fingerArray,
		Prime:             &prime,
		Seq:               req.Key.Seq,
	}
	raftOwner := &raftpb.Node{
		Ip:         req.Key.OwnerNode.Ip,
		Port:       req.Key.OwnerNode.Port,
		CommonName: req.Key.OwnerNode.CommonName,
		NodeId:     req.Key.OwnerNode.Uuid,
	}
	raftHolder := &raftpb.Node{
		Ip:         globals.ThisNode.IP,
		Port:       globals.ThisNode.Port,
		CommonName: globals.ThisNode.CommonName,
		NodeId:     globals.ThisNode.UUID,
	}

	err := globals.HeldKeyPieces.AddPiece(req.Key.Generation, req.Key.OwnerNode.Uuid, piece)
	if err != nil {
		return &pb.SendKeyPieceResponse{}, grpc.Errorf(codes.FailedPrecondition, "failed to save key piece to disk: %s", err)
	}
	Log.Info("Received KeyPiece from", req.Key.OwnerNode)
	if globals.RaftNetworkServer != nil && globals.RaftNetworkServer.State.Configuration.HasConfiguration() && req.AddElement {
		err := globals.RaftNetworkServer.RequestKeyStateUpdate(raftOwner, raftHolder, req.Key.Generation)
		if err != nil {
			if err == keyman.ErrGenerationDeprecated {
				return &pb.SendKeyPieceResponse{}, err
			} else {
				return &pb.SendKeyPieceResponse{}, grpc.Errorf(codes.FailedPrecondition, "failed to commit to Raft: %s", err)
			}
		}
	} else {
		return &pb.SendKeyPieceResponse{true}, nil
	}
	return &pb.SendKeyPieceResponse{}, nil
}
