package v1

import "time"

type State string

const (
	All     State = "all"
	Active  State = "active"
	Expired State = "expired"
)

type Metadata struct {
	CreationTimestamp time.Time     `json:"creationTimestamp"`
	TTL               time.Duration `json:"ttl"`
}

type Renew struct {
	TTL time.Duration `json:"ttl"`
}

type Count struct {
	Value int32 `json:"count"`
}

var Empty struct{}
