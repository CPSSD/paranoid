package server

import (
	pb "github.com/cpssd/paranoid/proto/fileserver"
	"golang.org/x/net/context"
	"strconv"
)

func (s *FileserverServer) ListServer(ctx context.Context, req *pb.ListServeRequest) (*pb.ListServeResponse, error) {
	var served []*pb.ServedFile
	for key := range FileMap {
		if FileMap[key].Uuid == req.Uuid {
			file := pb.ServedFile{
				FilePath:       FileMap[key].FilePath,
				FileHash:       key,
				AccessLimit:    FileMap[key].AccessLimit - FileMap[key].AccessAmmount,
				ExpirationTime: strconv.FormatInt(int64(FileMap[key].ExpirationTime.Minute()), 10)}
			served = append(served, &(file))
		}
	}
	return &pb.ListServeResponse{served}, nil
}
