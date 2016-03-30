package server

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/cpssd/paranoid/discovery-server/dnetserver"
	pb "github.com/cpssd/paranoid/proto/fileserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"strings"
	"time"
)

func (s *FileserverServer) ServeFile(ctx context.Context, req *pb.ServeRequest) (*pb.ServeResponse, error) {
	if req.Timeout <= 0 {
		req.Timeout = 1000
	}
	if req.Limit <= 0 {
		req.Limit = 1000
	}

	for i := 0; i < len(dnetserver.Nodes); i++ {
		if dnetserver.Nodes[i].Data.Uuid == req.Uuid {
			hasher := md5.New()
			fileUUID := req.FileName + req.Uuid
			filename := strings.Split(req.FileName, "/")
			hasher.Write([]byte(fileUUID))
			hash := hex.EncodeToString(hasher.Sum(nil))
			fileData := &FileCache{0,
				req.Limit,
				req.FileData,
				filename[len(filename)-1],
				time.Now().Add(time.Minute * time.Duration(req.Timeout))}
			FileMap[hash] = fileData
			return &pb.ServeResponse{hash, Port}, nil
		}
	}
	returnError := grpc.Errorf(codes.NotFound, "node was not found")
	return &pb.ServeResponse{"", ""}, returnError
}
