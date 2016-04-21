package pnetclient

import (
	"errors"
	"fmt"
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/keyman"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	raftpb "github.com/cpssd/paranoid/proto/raft"
	"golang.org/x/net/context"
)

func SendKeyPiece(uuid string, generation int64, piece *keyman.KeyPiece) error {
	node, err := globals.Nodes.GetNode(uuid)
	if err != nil {
		return errors.New("could not find node details")
	}

	conn, err := Dial(node)
	if err != nil {
		Log.Error("SendKeyPiece: failed to dial ", node)
		return fmt.Errorf("failed to dial: %s", node)
	}
	defer conn.Close()

	client := pb.NewParanoidNetworkClient(conn)

	thisNodeProto := &pb.Node{
		Ip:         globals.ThisNode.IP,
		Port:       globals.ThisNode.Port,
		CommonName: globals.ThisNode.CommonName,
		Uuid:       globals.ThisNode.UUID,
	}
	keyProto := &pb.KeyPiece{
		Data:              piece.Data,
		ParentFingerprint: piece.ParentFingerprint[:],
		Prime:             piece.Prime.Bytes(),
		Seq:               piece.Seq,
		Generation:        generation,
		OwnerNode:         thisNodeProto,
	}

	resp, err := client.SendKeyPiece(context.Background(), keyProto)
	if err != nil {
		Log.Error("Failed sending KeyPiece to", node, "Error:", err)
		return fmt.Errorf("Failed sending key piece to %s, Error: %s", node, err)
	}

	if resp.ClientMustCommit {
		raftThisNodeProto := &raftpb.Node{
			Ip:         globals.ThisNode.IP,
			Port:       globals.ThisNode.Port,
			CommonName: globals.ThisNode.CommonName,
			NodeId:     globals.ThisNode.UUID,
		}
		raftOwnerNode := &raftpb.Node{
			Ip:         keyProto.GetOwnerNode().Ip,
			Port:       keyProto.GetOwnerNode().Port,
			CommonName: keyProto.GetOwnerNode().CommonName,
			NodeId:     keyProto.GetOwnerNode().Uuid,
		}
		err := globals.RaftNetworkServer.RequestKeyStateUpdate(raftThisNodeProto, raftOwnerNode, generation)
		if err != nil {
			Log.Errorf("failed to commit to Raft: %s", err)
			return fmt.Errorf("failed to commit to Raft: %s", err)
		}
	}

	return nil
}
