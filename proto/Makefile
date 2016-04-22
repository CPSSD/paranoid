all: pnetwork dnetwork raftnetwork

pnetwork:
	protoc --go_out=plugins=grpc:. paranoidnetwork/paranoid.proto

dnetwork:
	protoc --go_out=plugins=grpc:. discoverynetwork/discovery.proto

raftnetwork:
	protoc --go_out=plugins=grpc:. raft/raft.proto

fileserver:
	protoc --go_out=plugins=grpc:. fileserver/fileserver.proto
