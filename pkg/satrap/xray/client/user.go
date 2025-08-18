package client

import (
	"context"

	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
)

func (c *Client) AddUser(ctx context.Context, tag string, user satrapv1.UserAccount) error {
	_, err := c.hsClient.AlterInbound(ctx, &command.AlterInboundRequest{
		Tag: tag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Email:   user.GetEmail(),
				Account: user.ToTypedMessage(),
			},
		}),
	})
	return handleXrayError(err)
}

func (c *Client) RemoveUser(ctx context.Context, tag, email string) error {
	_, err := c.hsClient.AlterInbound(ctx, &command.AlterInboundRequest{
		Tag: tag,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{
			Email: email,
		}),
	})
	return handleXrayError(err)
}
