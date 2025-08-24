package v1

import "time"

var Empty struct{}

type Count struct {
	Value int32 `json:"count"`
}

type Renew struct {
	TTL time.Duration `json:"ttl"`
}
