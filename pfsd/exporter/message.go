package exporter

type MessageType string

const (
	StateMessage      MessageType = "state"
	NodeChangeMessage             = "nodechange"
	RaftActionMessage             = "event"
)

func (m MessageType) String() string {
	switch m {
	case StateMessage:
		return "state"
	case NodeChangeMessage:
		return "nodechange"
	case RaftActionMessage:
		return "event"
	default:
		return ""
	}
}

type Message struct {
	Type MessageType `json:"type"`
	Data MessageData `json:"data"`
}

type MessageData struct {
	// Used for "status"
	Nodes []MessageNode `json:"nodes,omitempty"`
	// Used for "nodechange"
	Action string `json:"action,omitempty"`
	// Used for "nodechange"
	Node MessageNode `json:"node,omitempty"`
	// Used for "event"
	Event MessageEvent `json:"event,omitempty"`
}

type MessageNode struct {
	Uuid       string `json:"uuid"`
	CommonName string `json:"commonName"`
	State      string `json:"state"`
	Addr       string `json:"addr"`
}

type MessageEvent struct {
}
