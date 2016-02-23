package raft

import (
	"errors"
	"github.com/cpssd/paranoid/raft/raftlog"
)

//Configuration manages configuration information of a raft server
type Configuration struct {
	myNodeId                  string
	futureConfigurationActive bool

	currentConfiguration []Node
	currentNextIndex     []uint64
	currentMatchIndex    []uint64

	futureConfiguration []Node
	futureNextIndex     []uint64
	futureMatchIndex    []uint64
}

func (c *Configuration) GetNode(nodeID string) (Node, error) {
	for i := 0; i < len(c.currentConfiguration); i++ {
		if c.currentConfiguration[i].NodeID == nodeID {
			return c.currentConfiguration[i], nil
		}
	}

	for i := 0; i < len(c.futureConfiguration); i++ {
		if c.futureConfiguration[i].NodeID == nodeID {
			return c.futureConfiguration[i], nil
		}
	}

	return Node{}, errors.New("Node not found in configuration")
}

func (c *Configuration) NewFutureConfiguration(nodes []Node, lastLogIndex uint64) {
	c.futureConfigurationActive = true
	c.futureConfiguration = nodes
	c.futureNextIndex = make([]uint64, len(nodes))
	c.futureMatchIndex = make([]uint64, len(nodes))

	for i := 0; i < len(nodes); i++ {
		c.futureNextIndex[i] = lastLogIndex + 1
		c.futureMatchIndex[i] = 0
	}
}

func (c *Configuration) UpdateCurrentConfiguration(nodes []Node, lastLogIndex uint64) {
	if len(nodes) == len(c.futureConfiguration) {
		futureToCurrent := true
		for i := 0; i < len(nodes); i++ {
			if c.inFutureConfiguration(nodes[i].NodeID) == false {
				futureToCurrent = false
				break
			}
		}
		if futureToCurrent {
			c.FutureToCurrentConfiguration()
			return
		}
	}

	c.currentConfiguration = nodes
	c.currentNextIndex = make([]uint64, len(nodes))
	c.currentMatchIndex = make([]uint64, len(nodes))
	for i := 0; i < len(nodes); i++ {
		c.currentNextIndex[i] = lastLogIndex + 1
		c.currentMatchIndex[i] = 0
	}
}

func (c *Configuration) GetFutureConfigurationActive() bool {
	return c.futureConfigurationActive
}

func (c *Configuration) FutureToCurrentConfiguration() {
	c.futureConfigurationActive = false
	c.currentConfiguration = c.futureConfiguration
	c.currentNextIndex = c.futureNextIndex
	c.currentMatchIndex = c.futureMatchIndex

	c.futureConfiguration = []Node{}
	c.futureNextIndex = []uint64{}
	c.futureMatchIndex = []uint64{}
}

func (c *Configuration) inCurrentConfiguration(nodeID string) bool {
	for i := 0; i < len(c.currentConfiguration); i++ {
		if c.currentConfiguration[i].NodeID == nodeID {
			return true
		}
	}
	return false
}

func (c *Configuration) inFutureConfiguration(nodeID string) bool {
	for i := 0; i < len(c.futureConfiguration); i++ {
		if c.futureConfiguration[i].NodeID == nodeID {
			return true
		}
	}
	return false
}

func (c *Configuration) InConfiguration(nodeID string) bool {
	return c.inCurrentConfiguration(nodeID) || c.inFutureConfiguration(nodeID)
}

func (c *Configuration) MyConfigurationGood() bool {
	if c.InConfiguration(c.myNodeId) {
		if c.GetTotalPossibleVotes() > 1 {
			return true
		}
	}
	return false
}

func (c *Configuration) GetTotalPossibleVotes() int {
	votes := len(c.currentConfiguration)
	for i := 0; i < len(c.futureConfiguration); i++ {
		if c.inCurrentConfiguration(c.futureConfiguration[i].NodeID) == false {
			votes++
		}
	}
	return votes
}

func (c *Configuration) GetPeersList() []Node {
	var peers []Node
	for i := 0; i < len(c.currentConfiguration); i++ {
		if c.currentConfiguration[i].NodeID != c.myNodeId {
			peers = append(peers, c.currentConfiguration[i])
		}
	}
	for i := 0; i < len(c.futureConfiguration); i++ {
		if c.futureConfiguration[i].NodeID != c.myNodeId {
			if c.inCurrentConfiguration(c.futureConfiguration[i].NodeID) == false {
				peers = append(peers, c.futureConfiguration[i])
			}
		}
	}
	return peers
}

func getRequiredVotes(nodeCount int) int {
	if nodeCount == 0 {
		return 0
	}
	return (nodeCount / 2) + 1
}

//Check has a majority of votes have been received given a list of NodeIDs
func (c *Configuration) HasMajority(votesRecieved []string) bool {
	currentRequiredVotes := getRequiredVotes(len(c.currentConfiguration))
	count := 0
	for i := 0; i < len(votesRecieved); i++ {
		if c.inCurrentConfiguration(votesRecieved[i]) {
			count++
		}
	}
	if count < currentRequiredVotes {
		return false
	}

	if c.futureConfigurationActive {
		futureRequiredVotes := getRequiredVotes(len(c.futureConfiguration))
		count = 0
		for i := 0; i < len(votesRecieved); i++ {
			if c.inFutureConfiguration(votesRecieved[i]) {
				count++
			}
		}
		if count < futureRequiredVotes {
			return false
		}
	}
	return true
}

func (c *Configuration) ResetNodeIndexs(lastLogIndex uint64) {
	for i := 0; i < len(c.currentConfiguration); i++ {
		c.currentNextIndex[i] = lastLogIndex + 1
		c.currentMatchIndex[i] = 0
	}
	for i := 0; i < len(c.futureConfiguration); i++ {
		c.futureNextIndex[i] = lastLogIndex + 1
		c.futureMatchIndex[i] = 0
	}
}

func (c *Configuration) GetNextIndex(nodeID string) uint64 {
	for i := 0; i < len(c.currentConfiguration); i++ {
		if c.currentConfiguration[i].NodeID == nodeID {
			return c.currentNextIndex[i]
		}
	}
	for i := 0; i < len(c.futureConfiguration); i++ {
		if c.futureConfiguration[i].NodeID == nodeID {
			return c.futureNextIndex[i]
		}
	}
	Log.Fatal("Could not get NextIndex. Node not found")
	return 0
}

func (c *Configuration) GetMatchIndex(nodeID string) uint64 {
	for i := 0; i < len(c.currentConfiguration); i++ {
		if c.currentConfiguration[i].NodeID == nodeID {
			return c.currentMatchIndex[i]
		}
	}
	for i := 0; i < len(c.futureConfiguration); i++ {
		if c.futureConfiguration[i].NodeID == nodeID {
			return c.futureMatchIndex[i]
		}
	}
	Log.Fatal("Could not get MatchIndex. Node not found")
	return 0
}

func (c *Configuration) SetNextIndex(nodeID string, x uint64) {
	for i := 0; i < len(c.currentConfiguration); i++ {
		if c.currentConfiguration[i].NodeID == nodeID {
			c.currentNextIndex[i] = x
		}
	}
	for i := 0; i < len(c.futureConfiguration); i++ {
		if c.futureConfiguration[i].NodeID == nodeID {
			c.futureNextIndex[i] = x
		}
	}
}

func (c *Configuration) SetMatchIndex(nodeID string, x uint64) {
	for i := 0; i < len(c.currentConfiguration); i++ {
		if c.currentConfiguration[i].NodeID == nodeID {
			c.currentMatchIndex[i] = x
		}
	}
	for i := 0; i < len(c.futureConfiguration); i++ {
		if c.futureConfiguration[i].NodeID == nodeID {
			c.futureMatchIndex[i] = x
		}
	}
}

func (c *Configuration) CalculateNewCommitIndex(lastCommitIndex, term uint64, log *raftlog.RaftLog) uint64 {
	if log.GetMostRecentTerm() != term {
		return lastCommitIndex
	}

	currentMajority := getRequiredVotes(len(c.currentMatchIndex))
	futureMajority := getRequiredVotes(len(c.futureMatchIndex))
	newCommitIndex := lastCommitIndex

	for i := lastCommitIndex + 1; i <= log.GetMostRecentIndex(); i++ {
		logEntry, err := log.GetLogEntry(i)
		if err != nil {
			Log.Fatal("Unable to get log entry:", err)
		}
		if logEntry.Term == term {
			currentCount := 0
			if c.inCurrentConfiguration(c.myNodeId) {
				currentCount = 1
			}
			for j := 0; j < len(c.currentMatchIndex); j++ {
				if c.currentConfiguration[j].NodeID != c.myNodeId {
					if c.currentMatchIndex[j] >= i {
						currentCount++
					}
				}
			}
			if currentCount < currentMajority {
				return newCommitIndex
			}

			if c.futureConfigurationActive {
				futureCount := 0
				if c.inFutureConfiguration(c.myNodeId) {
					futureCount = 1
				}
				for j := 0; j < len(c.futureMatchIndex); j++ {
					if c.futureMatchIndex[j] >= i {
						futureCount++
					}
				}
				if futureCount < futureMajority {
					return newCommitIndex
				}

			}
			newCommitIndex = i
		}
	}
	return newCommitIndex
}

func newConfiguration(nodes []Node, nodeId string, lastLogIndex uint64) *Configuration {
	config := &Configuration{
		myNodeId:             nodeId,
		currentConfiguration: nodes,
		currentNextIndex:     make([]uint64, len(nodes)),
		currentMatchIndex:    make([]uint64, len(nodes)),
	}
	for i := 0; i < len(nodes); i++ {
		config.currentNextIndex[i] = lastLogIndex + 1
		config.currentMatchIndex[i] = 0
	}
	return config
}
