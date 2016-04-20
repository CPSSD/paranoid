package exporter

type MessageType int

const (
  NodeChangeMessage MessageType = iota
  RaftActionMessage
)

type Message struct {
  Type MessageType `json:"type"`
  Data MessageData `json:"data"`
}

type MessageData struct {
  // Used for "nodechange"
  Node MessageNode `json:"node,omitempty"`
  // Used for "status"
  Nodes []MessageNode `json:"nodes,omitempty"`
  // Used for Event
  Event MessageEvent `json:"event,omitempty"`
}

type MessageNode struct {
  Uuid string `json:"uuid"`
  CommonName string `json:"commonName"`
  State string `json:"state"`
}

type MessageEvent struct {

}
