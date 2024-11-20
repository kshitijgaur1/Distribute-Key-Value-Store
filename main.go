package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"distributed-kv-store/raft" // Replace with your correct local package path
	raftlib "github.com/hashicorp/raft"
)

func main() {
	// Clear old Raft state (for development purposes only)
	os.RemoveAll("./raft-data/node1")
	os.RemoveAll("./raft-data/node2")
	os.RemoveAll("./raft-data/node3")

	// Initialize Raft nodes
	node1 := raft.NewRaftNode("node1", "./raft-data/node1", "localhost:5001")
	// node2 := raft.NewRaftNode("node2", "./raft-data/node2", "localhost:5002")
	// node3 := raft.NewRaftNode("node3", "./raft-data/node3", "localhost:5003")

	// Bootstrap cluster on Node1
	configuration := raftlib.Configuration{
		Servers: []raftlib.Server{
			{ID: raftlib.ServerID("node1"), Address: raftlib.ServerAddress("localhost:5001")},
			{ID: raftlib.ServerID("node2"), Address: raftlib.ServerAddress("localhost:5002")},
			{ID: raftlib.ServerID("node3"), Address: raftlib.ServerAddress("localhost:5003")},
		},
	}

	log.Println("Bootstrapping cluster on Node1...")
	if err := node1.Raft.BootstrapCluster(configuration).Error(); err != nil {
		log.Fatalf("Failed to bootstrap cluster: %v", err)
	}
	log.Println("Cluster bootstrapped successfully.")

	// Delay to allow leader election
	time.Sleep(5 * time.Second)
	fmt.Println("Cluster initialized. Node1 is the leader.")

	// Check if Node1 is the leader
	if node1.Raft.State() == raftlib.Leader {
		// Simulate writing a key on Node1
		futureNode1 := node1.Raft.Apply([]byte("set:user123:Alice"), 10*time.Second)
		if err := futureNode1.Error(); err != nil {
			log.Fatalf("Node1 failed to apply log: %v", err)
		}
		log.Println("Log applied successfully on Node1")
	} else {
		log.Println("Node1 is not the leader. Aborting log application.")
	}
}
