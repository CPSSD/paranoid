package server

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/cpssd/paranoid/discovery-server/dnetserver"
	pb "github.com/cpssd/paranoid/proto/fileserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

func (s *FileserverServer) ServeFile(ctx context.Context, req *pb.ServeRequest) (*pb.ServeResponse, error) {
	if req.Timeout <= 0 {
		req.Timeout = 1000
	}
	if req.Limit <= 0 {
		req.Limit = 1000
	}

	for _, node := range dnetserver.Pools[req.Pool].Info.Nodes {
		if node.Uuid == req.Uuid {
			hasher := md5.New()
			fileUUID := req.FilePath + req.Uuid
			hasher.Write([]byte(fileUUID))
			hash := hex.EncodeToString(hasher.Sum(nil))
			fileData := &FileCache{
				req.Uuid,
				0, //The file hasnt been used yet
				req.Limit,
				req.FileData,
				req.FilePath,
				false,
				time.Now().Add(time.Minute * time.Duration(req.Timeout))}
			FileMap[hash] = fileData
			return &pb.ServeResponse{hash, Port}, nil
		}
	}
	returnError := grpc.Errorf(codes.NotFound, "node was not found")
	return &pb.ServeResponse{"", ""}, returnError
}
