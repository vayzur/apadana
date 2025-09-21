package v1

import (
	"encoding/json"
	"fmt"
	"time"

	metav1 "github.com/vayzur/apadana/pkg/apis/meta/v1"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/infra/conf"
	"github.com/xtls/xray-core/proxy/trojan"
	"github.com/xtls/xray-core/proxy/vless"
	"github.com/xtls/xray-core/proxy/vmess"
)

type Resource string

const (
	ResourceInbound Resource = "inbound"
	ResourceUser    Resource = "user"
)

type Count struct {
	Value uint32 `json:"count"`
}

type InboundCapacity struct {
	MaxUsers uint32 `json:"maxUsers"`
}

type InboundSpec struct {
	Capacity InboundCapacity          `json:"capacity"`
	Config   conf.InboundDetourConfig `json:"config"`
	TTL      time.Duration            `json:"ttl"`
}

type Inbound struct {
	Metadata metav1.ObjectMeta `json:"metadata"`
	Spec     InboundSpec       `json:"spec"`
}

type Account interface {
	ToTypedMessage() *serial.TypedMessage
}

type VlessAccount struct {
	ID   string
	Flow string
}

func (a VlessAccount) ToTypedMessage() *serial.TypedMessage {
	return serial.ToTypedMessage(&vless.Account{
		Id:   a.ID,
		Flow: a.Flow,
	})
}

type VmessAccount struct {
	ID string
}

func (a VmessAccount) ToTypedMessage() *serial.TypedMessage {
	return serial.ToTypedMessage(&vmess.Account{
		Id: a.ID,
	})
}

type TrojanAccount struct {
	Password string
}

func (a TrojanAccount) ToTypedMessage() *serial.TypedMessage {
	return serial.ToTypedMessage(&trojan.Account{
		Password: a.Password,
	})
}

type InboundUserSpec struct {
	Type       string          `json:"type"` // "vless", "vmess", "trojan"
	InboundTag string          `json:"inboundTag"`
	Email      string          `json:"email"`
	Account    json.RawMessage `json:"account"`
	TTL        time.Duration   `json:"ttl"`
}

type InboundUser struct {
	Metadata metav1.ObjectMeta `json:"metadata"`
	Spec     InboundUserSpec   `json:"spec"`
}

func (u *InboundUser) ToAccount() (Account, error) {
	switch u.Spec.Type {
	case "vless":
		var v VlessAccount
		if err := json.Unmarshal(u.Spec.Account, &v); err != nil {
			return nil, err
		}
		return &v, nil
	case "vmess":
		var v VmessAccount
		if err := json.Unmarshal(u.Spec.Account, &v); err != nil {
			return nil, err
		}
		return &v, nil
	case "trojan":
		var t TrojanAccount
		if err := json.Unmarshal(u.Spec.Account, &t); err != nil {
			return nil, err
		}
		return &t, nil
	default:
		return nil, fmt.Errorf("unknown protocol: %s", u.Spec.Type)
	}
}
