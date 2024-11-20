package hashing

import (
	"hash/crc32"
	"sort"
	// "strconv"
)

// HashRing represents a consistent hashing ring
type HashRing struct {
	nodes      map[uint32]string // Maps hash to node address
	sortedKeys []uint32          // Sorted hash values for the ring
}

// NewHashRing initializes a new HashRing
func NewHashRing() *HashRing {
	return &HashRing{
		nodes:      make(map[uint32]string),
		sortedKeys: []uint32{},
	}
}

// AddNode adds a node to the hash ring
func (hr *HashRing) AddNode(node string) {
	hash := crc32.ChecksumIEEE([]byte(node))
	hr.nodes[hash] = node
	hr.sortedKeys = append(hr.sortedKeys, hash)
	sort.Slice(hr.sortedKeys, func(i, j int) bool { return hr.sortedKeys[i] < hr.sortedKeys[j] })
}

// RemoveNode removes a node from the hash ring
func (hr *HashRing) RemoveNode(node string) {
	hash := crc32.ChecksumIEEE([]byte(node))
	delete(hr.nodes, hash)
	hr.rebuildSortedKeys()
}

// GetNode finds the appropriate node for a given key
func (hr *HashRing) GetNode(key string) string {
	hash := crc32.ChecksumIEEE([]byte(key))
	for _, nodeHash := range hr.sortedKeys {
		if hash <= nodeHash {
			return hr.nodes[nodeHash]
		}
	}
	return hr.nodes[hr.sortedKeys[0]] // Wrap around to the first node
}

// rebuildSortedKeys rebuilds the sorted hash keys
func (hr *HashRing) rebuildSortedKeys() {
	hr.sortedKeys = []uint32{}
	for hash := range hr.nodes {
		hr.sortedKeys = append(hr.sortedKeys, hash)
	}
	sort.Slice(hr.sortedKeys, func(i, j int) bool { return hr.sortedKeys[i] < hr.sortedKeys[j] })
}

// SortedNodes returns the sorted keys (hashes) of nodes
func (hr *HashRing) SortedNodes() []uint32 {
	return hr.sortedKeys
}

// NodeAtHash retrieves a node by its hash
func (hr *HashRing) NodeAtHash(hash uint32) string {
	return hr.nodes[hash]
}

// IndexOf finds the index of a node in the sorted hash list
func (hr *HashRing) IndexOf(node string) int {
	hash := crc32.ChecksumIEEE([]byte(node))
	for i, h := range hr.sortedKeys {
		if h == hash {
			return i
		}
	}
	return -1
}

// HashKey hashes a key to an integer
func (hr *HashRing) HashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}
