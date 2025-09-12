package v1

import (
	"github.com/xtls/xray-core/infra/conf"
)

type InboundCapacity struct {
	MaxUsers uint32 `json:"maxUsers"`
}

type InboundSpec struct {
	Capacity InboundCapacity          `json:"capacity"`
	Config   conf.InboundDetourConfig `json:"config"`
}

type Inbound struct {
	Metadata Metadata    `json:"metadata"`
	Spec     InboundSpec `json:"spec"`
}
