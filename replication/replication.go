package replication

import (
	"distributed-kv-store/hashing"
	"fmt"
)

type ReplicationManager struct {
	HashRing          *hashing.HashRing
	ReplicationFactor int
}

// NewReplicationManager initializes a replication manager
func NewReplicationManager(hashRing *hashing.HashRing, replicationFactor int) *ReplicationManager {
	return &ReplicationManager{
		HashRing:          hashRing,
		ReplicationFactor: replicationFactor,
	}
}

// GetNodesForKey retrieves the primary and replica nodes for a key
func (rm *ReplicationManager) GetNodesForKey(key string) []string {
	allNodes := rm.HashRing.SortedNodes()
	hash := rm.HashRing.HashKey(key)
	result := []string{}

	// Find the primary node
	for _, nodeHash := range allNodes {
		if hash <= nodeHash {
			result = append(result, rm.HashRing.NodeAtHash(nodeHash))
			break
		}
	}
	if len(result) == 0 {
		// Wrap around to the first node
		result = append(result, rm.HashRing.NodeAtHash(allNodes[0]))
	}

	// Add replica nodes
	startIndex := rm.HashRing.IndexOf(result[0])
	for i := 1; i < rm.ReplicationFactor; i++ {
		nextIndex := (startIndex + i) % len(allNodes)
		result = append(result, rm.HashRing.NodeAtHash(allNodes[nextIndex]))
	}

	return result
}

// PrintReplication shows the nodes responsible for a key
func (rm *ReplicationManager) PrintReplication(key string) {
	nodes := rm.GetNodesForKey(key)
	fmt.Printf("Key '%s' is stored in nodes: %v\n", key, nodes)
}
