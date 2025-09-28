package v1

import (
	"time"

	metav1 "github.com/vayzur/apadana/pkg/apis/meta/v1"
)

const (
	LabelHostname = "hostname"
	LabelOS       = "os"
	LabelArch     = "arch"
	LabelCountry  = "country"
	LabelRegion   = "region"
	LabelProvider = "provider"
)

type NodeAddressType string

const (
	InternalAddress NodeAddressType = "InternalAddress"
	ExternalAddress NodeAddressType = "ExternalAddress"
)

type NodeAddress struct {
	Type    NodeAddressType `json:"type"`
	Address string          `json:"address"`
}

type NodeCapacity struct {
	MaxInbounds uint32 `json:"maxInbounds"`
}

type NodeStatus struct {
	Capacity          NodeCapacity  `json:"capacity"`
	Addresses         []NodeAddress `json:"addresses"`
	Ready             bool          `json:"ready"`
	LastHeartbeatTime time.Time     `json:"lastHeartbeatTime"`
}

type Node struct {
	Metadata metav1.ObjectMeta `json:"metadata"`
	Status   NodeStatus        `json:"status"`
}

func GetPreferredAddress(addresses []NodeAddress, addressType NodeAddressType) string {
	for _, addr := range addresses {
		if addr.Type == addressType {
			return addr.Address
		}
	}
	return addresses[0].Address
}
