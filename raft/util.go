package raft

import(
  pb "github.com/cpssd/paranoid/proto/raft"
  "github.com/cpssd/paranoid/pfsd/exporter"
)

func convertProtoDetailedNodeToExportNode(nodes []*pb.LeaderData_Data_DetailedNode) []exporter.MessageNode {
  var res []exporter.MessageNode

  for i := 0; i < len(nodes); i++ {
    res = append(res, exporter.MessageNode{
      CommonName: nodes[i].CommonName,
      Addr: nodes[i].Addr,
      Uuid: nodes[i].Uuid,
      State: nodes[i].State,
    })
  }
  return res
}
