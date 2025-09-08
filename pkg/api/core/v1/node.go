package v1

import (
	"fmt"
	"time"
)

type NodeAddressType string

const (
	InternalAddress NodeAddressType = "InternalAddress"
	ExternalAddress NodeAddressType = "ExternalAddress"
)

type NodeMetadata struct {
	Name              string            `json:"name"`
	ID                string            `json:"id"`
	CreationTimestamp time.Time         `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
}

type NodeSpec struct {
	Token string `json:"token"`
}

type NodeAddress struct {
	Type NodeAddressType `json:"type"`
	Host string          `json:"host"`
}

type NodeCapacity struct {
	MaxInbounds uint32 `json:"maxInbounds"`
}

type NodeConnectionConfig struct {
	Scheme string `json:"scheme"`
	Port   uint16 `json:"port"`
}

type NodeStatus struct {
	Addresses         []NodeAddress        `json:"addresses"`
	Capacity          NodeCapacity         `json:"capacity"`
	Ready             bool                 `json:"ready"`
	LastHeartbeatTime time.Time            `json:"lastHeartbeatTime"`
	ConnectionConfig  NodeConnectionConfig `json:"connectionConfig"`
}

type Node struct {
	Metadata NodeMetadata `json:"metadata"`
	Spec     NodeSpec     `json:"spec"`
	Status   NodeStatus   `json:"status"`
}

func getPreferredAddress(addresses []NodeAddress) string {
	for _, addr := range addresses {
		if addr.Type == InternalAddress {
			return addr.Host
		}
	}
	return addresses[0].Host
}

func (n *Node) URL(path string) string {
	host := getPreferredAddress(n.Status.Addresses)
	return fmt.Sprintf("%s://%s:%d%s", n.Status.ConnectionConfig.Scheme, host, n.Status.ConnectionConfig.Port, path)
}
