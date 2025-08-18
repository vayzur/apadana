package v1

import "time"

type NodeMetadata struct {
	Name              string    `json:"name"`
	ID                string    `json:"id"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}

type NodeStatus struct {
	Status            bool      `json:"status"`
	LastHeartbeatTime time.Time `json:"lastHeartbeatTime"`
}

type Node struct {
	Metadata NodeMetadata `json:"metadata"`
	Status   NodeStatus   `json:"status"`
	Address  string       `json:"address"`
	Token    string       `json:"token"`
}
