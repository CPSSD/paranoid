package pnetserver

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	raftpb "github.com/cpssd/paranoid/proto/raft"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// NewGeneration recieves requests from nodes asking to create a new KeyPiece
// generation in preparation for joining the cluster.
func (s *ParanoidServer) NewGeneration(ctx context.Context, req *pb.NewGenerationRequest) (*pb.NewGenerationResponse, error) {
	if req.PoolPassword == "" {
		if len(globals.PoolPasswordHash) != 0 {
			return &pb.NewGenerationResponse{}, grpc.Errorf(codes.InvalidArgument,
				"cluster is password protected but no password was given")
		}
	} else {
		err := bcrypt.CompareHashAndPassword(globals.PoolPasswordHash,
			append(globals.PoolPasswordSalt, []byte(req.PoolPassword)...))
		if err != nil {
			return &pb.NewGenerationResponse{}, grpc.Errorf(codes.InvalidArgument,
				"unable to request new generation: password error:", err)
		}
	}

	raftReqNode := &raftpb.Node{
		Ip:         req.GetRequestingNode().Ip,
		Port:       req.GetRequestingNode().Port,
		CommonName: req.GetRequestingNode().CommonName,
		NodeId:     req.GetRequestingNode().Uuid,
	}
	generationNumber, peers, err := globals.RaftNetworkServer.RequestNewGeneration(raftReqNode)
	if err != nil {
		return &pb.NewGenerationResponse{}, grpc.Errorf(codes.Unknown, "unable to create new generation")
	}
	// We need to convert between the two different Node types.
	paranoidPeers := make([]*pb.Node, len(peers))
	for i, v := range peers {
		paranoidPeers[i] = &pb.Node{
			Ip:         v.Ip,
			Port:       v.Port,
			CommonName: v.CommonName,
			Uuid:       v.NodeId,
		}
	}
	return &pb.NewGenerationResponse{
		GenerationNumber: int64(generationNumber),
		Peers:            paranoidPeers,
	}, nil
}
