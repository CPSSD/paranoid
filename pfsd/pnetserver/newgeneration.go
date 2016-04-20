package pnetserver

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
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

	generationNumber, peers, err := globals.RaftNetworkServer.RequestNewGeneration(req.GetRequestingNode().Uuid)
	if err != nil {
		return &pb.NewGenerationResponse{}, grpc.Errorf(codes.Unknown, "unable to create new generation")
	}
	// We need to convert between the two different Node types.
	paranoidPeers := make([]*pb.Node, len(peers))
	for i, v := range peers {
		globalNode, err := globals.Nodes.GetNode(v)
		if err != nil {
			Log.Errorf("Could not look up node %s: %s", v, err)
		}
		paranoidPeers[i] = &pb.Node{
			Ip:         globalNode.IP,
			Port:       globalNode.Port,
			CommonName: globalNode.CommonName,
			Uuid:       globalNode.UUID,
		}
	}
	return &pb.NewGenerationResponse{
		GenerationNumber: int64(generationNumber),
		Peers:            paranoidPeers,
	}, nil
}
