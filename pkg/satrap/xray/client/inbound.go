package client

import (
	"context"

	satrapv1 "github.com/vayzur/apadana/pkg/apis/satrap/v1"
	"github.com/vayzur/apadana/pkg/errs"
	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/infra/conf"
)

func (c *Client) ListInbounds(ctx context.Context) (map[string]struct{}, error) {
	reqCtx, cancel := c.withTimeout(ctx)
	defer cancel()
	resp, err := c.hsClient.ListInbounds(reqCtx, &command.ListInboundsRequest{
		IsOnlyTags: true,
	})
	if err != nil {
		return nil, errs.New(errs.KindInternal, errs.ReasonUnknown, "list inbounds failed", nil, err)
	}

	inbs := resp.GetInbounds()

	inbounds := make(map[string]struct{}, len(inbs))
	for _, inbound := range inbs {
		inbounds[inbound.Tag] = struct{}{}
	}

	// remove the "api" tag in one operation
	delete(inbounds, "api")

	return inbounds, nil
}

func (c *Client) AddInbound(ctx context.Context, conf *conf.InboundDetourConfig) error {
	config, err := conf.Build()
	if err != nil {
		return errs.New(errs.KindInvalid, errs.ReasonUnknown, "inbound config build failed", nil, err)
	}

	reqCtx, cancel := c.withTimeout(ctx)
	defer cancel()

	_, err = c.hsClient.AddInbound(reqCtx, &command.AddInboundRequest{
		Inbound: config,
	})
	return errs.HandleXrayError(err, satrapv1.ResourceInbound)
}

func (c *Client) RemoveInbound(ctx context.Context, tag string) error {
	reqCtx, cancel := c.withTimeout(ctx)
	defer cancel()
	_, err := c.hsClient.RemoveInbound(reqCtx, &command.RemoveInboundRequest{
		Tag: tag,
	})
	return errs.HandleXrayError(err, satrapv1.ResourceInbound)
}

func (c *Client) AddUser(ctx context.Context, tag, email string, account satrapv1.Account) error {
	reqCtx, cancel := c.withTimeout(ctx)
	defer cancel()
	_, err := c.hsClient.AlterInbound(reqCtx, &command.AlterInboundRequest{
		Tag: tag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Email:   email,
				Account: account.ToTypedMessage(),
			},
		}),
	})
	return errs.HandleXrayError(err, satrapv1.ResourceUser)
}

func (c *Client) RemoveUser(ctx context.Context, tag, email string) error {
	reqCtx, cancel := c.withTimeout(ctx)
	defer cancel()
	_, err := c.hsClient.AlterInbound(reqCtx, &command.AlterInboundRequest{
		Tag: tag,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{
			Email: email,
		}),
	})
	return errs.HandleXrayError(err, satrapv1.ResourceUser)
}

func (c *Client) ListUsers(ctx context.Context, tag string) (map[string]struct{}, error) {
	reqCtx, cancel := c.withTimeout(ctx)
	defer cancel()
	resp, err := c.hsClient.GetInboundUsers(reqCtx, &command.GetInboundUserRequest{
		Tag: tag,
	})
	if err != nil {
		return nil, errs.New(errs.KindInternal, errs.ReasonUnknown, "list users failed", nil, err)
	}

	u := resp.GetUsers()

	users := make(map[string]struct{}, len(u))
	for _, user := range u {
		users[user.Email] = struct{}{}
	}

	return users, nil
}
