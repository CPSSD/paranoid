all: pnetwork dnetwork

pnetwork:
	protoc --go_out=plugins=grpc:. paranoidnetwork/paranoid.proto

dnetwork:
	protoc --go_out=plugins=grpc:. discoverynetwork/discovery.proto
