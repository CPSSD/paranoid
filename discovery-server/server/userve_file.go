package server

import (
	"github.com/cpssd/paranoid/discovery-server/dnetserver"
	pb "github.com/cpssd/paranoid/proto/fileserver"
	"golang.org/x/net/context"
)

func (s *FileserverServer) UnServeFile(ctx context.Context, req *pb.UnServeRequest) (*pb.ServeResponse, error) {
	for i := 0; i < len(dnetserver.Nodes); i++ {
		if dnetserver.Nodes[i].Data.Uuid == req.Uuid {
			delete(FileMap, req.TileHash)
		}
	}
	return &pb.ServeResponse{"File Successfully Removed", ""}, nil
}
