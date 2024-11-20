package raft

import (
	"log"
	"os"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

type RaftNode struct {
	Raft *raft.Raft
}

func NewRaftNode(nodeID, dataDir, bindAddress string) *RaftNode {
	// Raft configuration
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)

	// Create FSM
	fsm := NewFSM()

	// Log store
	logStore, err := raftboltdb.NewBoltStore(dataDir + "/raft-log.bolt")
	if err != nil {
		log.Fatalf("Failed to create log store: %v", err)
	}

	// Stable store
	stableStore, err := raftboltdb.NewBoltStore(dataDir + "/raft-stable.bolt")
	if err != nil {
		log.Fatalf("Failed to create stable store: %v", err)
	}

	// Snapshot store
	snapshots, err := raft.NewFileSnapshotStore(dataDir, 2, os.Stderr)
	if err != nil {
		log.Fatalf("Failed to create snapshot store: %v", err)
	}

	// Transport
	transport, err := raft.NewTCPTransport(bindAddress, nil, 3, raft.DefaultTimeout, os.Stderr)
	if err != nil {
		log.Fatalf("Failed to create transport: %v", err)
	}

	// Create Raft system
	r, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshots, transport)
	if err != nil {
		log.Fatalf("Failed to create Raft node: %v", err)
	}

	return &RaftNode{Raft: r}
}
