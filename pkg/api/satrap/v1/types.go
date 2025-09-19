package v1

import "time"

type Resource string

const (
	ResourceInbound  Resource = "inbound"
	ResourceUser     Resource = "user"
	ResourceOutbound Resource = "outbound"
)

type Metadata struct {
	CreationTimestamp time.Time         `json:"creationTimestamp"`
	TTL               time.Duration     `json:"ttl"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
}

type Renew struct {
	TTL time.Duration `json:"ttl"`
}

type Count struct {
	Value uint32 `json:"count"`
}
