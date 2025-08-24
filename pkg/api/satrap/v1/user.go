package v1

import (
	"encoding/json"
	"fmt"

	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/proxy/trojan"
	"github.com/xtls/xray-core/proxy/vless"
	"github.com/xtls/xray-core/proxy/vmess"
)

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

type InboundUser struct {
	Metadata Metadata        `json:"metadata"`
	Type     string          `json:"type"` // "vless", "vmess", "trojan"
	Email    string          `json:"email"`
	Account  json.RawMessage `json:"account"`
}

func (u *InboundUser) ToAccount() (Account, error) {
	switch u.Type {
	case "vless":
		var v VlessAccount
		if err := json.Unmarshal(u.Account, &v); err != nil {
			return nil, err
		}
		return &v, nil
	case "vmess":
		var v VmessAccount
		if err := json.Unmarshal(u.Account, &v); err != nil {
			return nil, err
		}
		return &v, nil
	case "trojan":
		var t TrojanAccount
		if err := json.Unmarshal(u.Account, &t); err != nil {
			return nil, err
		}
		return &t, nil
	default:
		return nil, fmt.Errorf("unknown protocol: %s", u.Type)
	}
}
