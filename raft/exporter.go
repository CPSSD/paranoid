package raft

import (
  "time"
)

const (
  LEADER_CHECK_INTERVAL time.Duration = 100 * time.Millisecond
)

func (s *RaftNetworkServer) sendLeaderDataRequest() {
	defer s.Wait.Done()

	leaderCheckTimer := time.NewTimer(LEADER_CHECK_INTERVAL)

checkLeaderLoop:
	for {
		select {
		case <-leaderCheckTimer.C:
			if s.getLeader() == nil {
				Log.Warn("Leader not yet elected")
				continue
			} else {
				Log.Info("Leader elected")
				leaderCheckTimer.Stop()
				break checkLeaderLoop
			}
		}
	}

dialLeader:
	leader := s.getLeader()
	conn, err := s.Dial(s.getLeader(), SEND_ENTRY_TIMEOUT)
	if err != nil {
		Log.Error("Unable to dial leader")
	}
	client := pb.NewRaftNetworkClient(conn)
	stream, err := client.RequestLeaderData(context.Background(), &pb.LeaderDataRequest{})
	if err != nil {
		Log.Error("Unable to request user data")
	}

	for {
		select {
		case <-s.Quit:
		default:
			// Check is the leader we are dialing still the leader
			if(leader != s.getLeader()){
				goto dialLeader
			}

			data, err := stream.Recv()
			if err != nil {
				Log.Error("Unable to get data:", err)
			}

			// Get the message from protobuf and convert it to export message
			var messageType exporter.MessageType
			var messageData exporter.MessageData

			switch data.Type {
			case pb.LeaderData_State:
				messageType = exporter.StateMessage
				messageData = exporter.MessageData{
					Nodes: convertProtoDetailedNodeToExportNode(data.Data.GetNodes()),
				}
			case pb.LeaderData_NodeChange:
				messageType = exporter.NodeChangeMessage
				messageData = exporter.MessageData{
					Action: data.Data.Action,
					Node: exporter.MessageNode{
						CommonName: data.Data.Node.CommonName,
			      Addr: data.Data.Node.Addr,
			      Uuid: data.Data.Node.Uuid,
			      State: data.Data.Node.State,
					},
				}
			case pb.LeaderData_Event:
				messageType = exporter.RaftEventMessage
				messageData = exporter.MessageData{
					Event: exporter.MessageEvent{
						Source: data.Data.Event.Source,
						Target: data.Data.Event.Target,
						Details: data.Data.Event.Details,
					},
				}
			}

			msg := exporter.Message{
				Type: messageType,
				Data: messageData,
			}

			// Send the export message
			exporter.Send(msg)
		}
	}
}

//sendLeaderLogEntry forwards a client request to the leader
func (s *RaftNetworkServer) sendLeaderLogEntry(entry *pb.Entry) error {
	sendLogTimeout := time.After(LEADER_REQUEST_TIMEOUT)
	lastError := errors.New("timeout before client to leader request was attempted")
	for {
		select {
		case <-sendLogTimeout:
			return lastError
		default:
			leaderNode := s.getLeader()
			if leaderNode == nil {
				lastError = errors.New("Unable to find leader")
				continue
			}

			conn, err := s.Dial(leaderNode, SEND_ENTRY_TIMEOUT)
			if err != nil {
				lastError = err
				continue
			}
			defer conn.Close()

			if err == nil {
				client := pb.NewRaftNetworkClient(conn)
				_, err := client.ClientToLeaderRequest(context.Background(), &pb.EntryRequest{s.State.NodeId, entry})
				if err == nil {
					return err
				}
				lastError = err
			}
		}
	}
}

func updateExporterState(nodes []Node){
  var detailedNodes []*pb.LeaderData_Data_DetailedNode
  for i := 0; i < len(nodes); i++ {
    detailedNodes = append(detailedNodes, &pb.LeaderData_Data_DetailedNode{
      Uuid: nodes[i].NodeID,
      CommonName: nodes[i].CommonName,
      State: "unknown",
      Addr: nodes[i].IP+":"+nodes[i].Port,
    })
  }

  // Send the status to listening channel
  exportedChangeList <- pb.LeaderData{
    Type: pb.LeaderData_State,
    Data: &pb.LeaderData_Data{
      Nodes: detailedNodes,
    },
  }
}
