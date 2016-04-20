package pnetclient

import (
	"github.com/cpssd/paranoid/pfsd/globals"
	"github.com/cpssd/paranoid/pfsd/keyman"
	pb "github.com/cpssd/paranoid/proto/paranoidnetwork"
	raftpb "github.com/cpssd/paranoid/proto/raft"
	"golang.org/x/net/context"
)

func SendKeyPiece(piece *keyman.KeyPiece) {
	// We use piece.Seq-2 because Seq1 goes to the current node.
	// We only distribute Seq2 onwards.
	node := globals.Nodes.GetAll()[piece.Seq-2]
	conn, err := Dial(node)
	if err != nil {
		Log.Error("SendKeyPiece: failed to dial ", node)
		return
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
		OwnerNode:         thisNodeProto,
	}
	resp, err := client.SendKeyPiece(context.Background(), keyProto)
	if err != nil {
		Log.Error("Failed sending KeyPiece to", node, "Error:", err)
		return
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
		err := globals.RaftNetworkServer.RequestKeyStateUpdate(raftThisNodeProto, raftOwnerNode,
			int64(keyman.StateMachine.CurrentGeneration+1))
		if err != nil {
			Log.Errorf("failed to commit to Raft: %s", err)
			return
		}
	}
}
