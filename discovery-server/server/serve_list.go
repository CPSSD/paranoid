package server

import (
	pb "github.com/cpssd/paranoid/proto/fileserver"
	"golang.org/x/net/context"
)

func (s *FileserverServer) ListServer(ctx context.Context, req *pb.ServeRequest) (*pb.ListServeResponse, error) {
	// var served []pb.ServedFiles
	// for key := range FileMap {
	// 	if FileMap[key].Uuid == req.Uuid {
	// 		file := pb.ServedFiles{
	// 			FilePath:       FileMap[key].FilePath,
	// 			FileHash:       key,
	// 			AccessLimit:    FileMap[key].AccessAmmount,
	// 			ExpirationTime: "FileMap[key].ExpirationTime"}
	// 		served = append(served, &(file))
	// 	}
	// }
	return nil, nil
}
