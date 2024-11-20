package raft

import (
	"io"        // Required for the FSM Restore method
	"log"       // For logging
	"os"        // For file and directory operations
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"         // Raft protocol
	"github.com/hashicorp/raft-boltdb" // BoltDB for Raft storage
)

// RaftNode represents a single node in the Raft cluster
type RaftNode struct {
	Raft *raft.Raft
}

// NewRaftNode initializes a Raft node
func NewRaftNode(nodeID, raftDir, bindAddr string) *RaftNode {
	// Create Raft configuration
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)

	// Setup storage directories
	if err := os.MkdirAll(raftDir, 0755); err != nil {
		log.Fatalf("Failed to create directory %s: %v", raftDir, err)
	}

	// Log store
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft-log.bolt"))
	if err != nil {
		log.Fatalf("Failed to create log store: %v", err)
	}

	// Stable store
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft-stable.bolt"))
	if err != nil {
		log.Fatalf("Failed to create stable store: %v", err)
	}

	// Snapshot store
	snapshotStore := raft.NewDiscardSnapshotStore()

	// Transport layer
	transport, err := raft.NewTCPTransport(bindAddr, nil, 3, 10*time.Second, os.Stderr)
	if err != nil {
		log.Fatalf("Failed to create TCP transport: %v", err)
	}

	// Finite State Machine (FSM)
	fsm := &FSM{}

	// Initialize Raft
	r, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		log.Fatalf("Failed to initialize Raft: %v", err)
	}

	return &RaftNode{Raft: r}
}

// FSM is the finite state machine for Raft
type FSM struct{}

// Apply applies a log entry to the FSM
// Apply applies a log entry to the FSM
func (f *FSM) Apply(entry *raft.Log) interface{} {
	log.Printf("Applied log entry: %s", string(entry.Data)) // Use 'entry.Data' for log content
	return nil
}

// Snapshot creates a snapshot of the FSM
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &FSMSnapshot{}, nil
}

// Restore restores the FSM from a snapshot
func (f *FSM) Restore(snapshot io.ReadCloser) error {
	return nil
}

// FSMSnapshot represents a snapshot of the FSM
type FSMSnapshot struct{}

// Persist saves the snapshot
func (f *FSMSnapshot) Persist(sink raft.SnapshotSink) error {
	return nil
}

// Release releases the snapshot resources
func (f *FSMSnapshot) Release() {}
