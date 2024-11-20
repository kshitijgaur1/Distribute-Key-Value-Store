package raft

import (
	"log"              // Use the Go standard log package
	"github.com/hashicorp/raft"
	"sync"
)

// FSM represents the finite state machine that holds the key-value store.
type FSM struct {
	data map[string]string // In-memory key-value store
	mu   sync.Mutex        // Mutex to protect access to the data
}

// NewFSM creates a new FSM instance with an empty key-value store.
func NewFSM() *FSM {
	return &FSM{
		data: make(map[string]string),
	}
}

// Apply is called by the Raft log to apply a log entry to the FSM.
func (fsm *FSM) Apply(log *raft.Log) interface{} {
	// Log the applied log entry using the standard log package (not raft.Log)
	log.Printf("Applying log entry: Index: %d, Term: %d, Type: %s", log.Index, log.Term, log.Type.String())

	// Ensure thread-safe access to the key-value store
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	// Check the log type and apply accordingly
	if log.Type == raft.LogCommand {
		// Deserialize the command from the log entry
		command := string(log.Data)
		parts := splitCommand(command) // A helper function to parse command

		// Assuming command format "SET key value"
		if len(parts) == 3 && parts[0] == "SET" {
			key := parts[1]
			value := parts[2]
			fsm.data[key] = value
			log.Printf("SET command applied: %s = %s", key, value)
		} else {
			log.Printf("Invalid command format: %s", command)
		}
	} else {
		log.Printf("Unknown log type: %s", log.Type.String())
	}
	return nil
}

// Snapshot returns the current snapshot of the FSM (key-value store).
func (fsm *FSM) Snapshot() ([]byte, error) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	// Serialize the key-value store data to create a snapshot
	var snapshot []byte
	for key, value := range fsm.data {
		snapshot = append(snapshot, []byte(key+"="+value+"\n")...)
	}
	log.Printf("Snapshot taken with %d entries", len(fsm.data))
	return snapshot, nil
}

// Restore restores the FSM state from a snapshot.
func (fsm *FSM) Restore(snapshot []byte) error {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	// Deserialize the snapshot back into the FSM state
	fsm.data = make(map[string]string)
	for _, line := range string(snapshot) {
		parts := splitCommand(line)
		if len(parts) == 2 {
			fsm.data[parts[0]] = parts[1]
		}
	}

	log.Printf("FSM restored from snapshot with %d entries", len(fsm.data))
	return nil
}

// Helper function to split the SET command into parts.
func splitCommand(command string) []string {
	// Split the command by spaces
	// Example: "SET key value" => ["SET", "key", "value"]
	parts := []string{}
	words := []byte{}
	for _, c := range command {
		if c == ' ' {
			parts = append(parts, string(words))
			words = []byte{}
		} else {
			words = append(words, byte(c))
		}
	}
	parts = append(parts, string(words)) // Append the last part
	return parts
}
