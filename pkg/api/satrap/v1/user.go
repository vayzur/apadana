package v1

import (
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/proxy/trojan"
	"github.com/xtls/xray-core/proxy/vless"
	"github.com/xtls/xray-core/proxy/vmess"
)

type UserAccount interface {
	ToTypedMessage() *serial.TypedMessage
	GetEmail() string
}

type BaseUser struct {
	Email string
}

type VlessUser struct {
	BaseUser
	ID   string
	Flow string
}

func (u VlessUser) ToTypedMessage() *serial.TypedMessage {
	return serial.ToTypedMessage(&vless.Account{
		Id:   u.ID,
		Flow: u.Flow,
	})
}

func (u VlessUser) GetEmail() string { return u.Email }

type VmessUser struct {
	BaseUser
	ID string
}

func (u VmessUser) ToTypedMessage() *serial.TypedMessage {
	return serial.ToTypedMessage(&vmess.Account{
		Id: u.ID,
	})
}

func (u VmessUser) GetEmail() string { return u.Email }

type TrojanUser struct {
	BaseUser
	Password string
}

func (u TrojanUser) ToTypedMessage() *serial.TypedMessage {
	return serial.ToTypedMessage(&trojan.Account{
		Password: u.Password,
	})
}

func (u TrojanUser) GetEmail() string { return u.Email }
