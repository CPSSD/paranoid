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

	err := globals.HeldKeyPieces.AddPiece(req.Generation, req.OwnerNode.Uuid, piece)
	if err != nil {
		return &pb.SendKeyPieceResponse{}, grpc.Errorf(codes.FailedPrecondition, "failed to save key piece to disk: %s", err)
	}
	Log.Info("Received KeyPiece from", req.OwnerNode)
	if globals.RaftNetworkServer != nil {
		err := globals.RaftNetworkServer.RequestKeyStateUpdate(raftOwner, raftHolder, req.Generation)
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
