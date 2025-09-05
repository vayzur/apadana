package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
	"github.com/vayzur/apadana/pkg/errs"
	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/infra/conf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) ListInbounds(ctx context.Context) (map[string]struct{}, error) {
	req := &command.ListInboundsRequest{IsOnlyTags: true}
	inbounds, err := c.hsClient.ListInbounds(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("list inbounds failed: %w", err)
	}

	inbs := make(map[string]struct{}, len(inbounds.Inbounds))
	for _, inb := range inbounds.Inbounds {
		inbs[inb.Tag] = satrapv1.Empty
	}

	// remove the "api" tag in one operation
	delete(inbs, "api")

	return inbs, nil
}

func (c *Client) AddInbound(ctx context.Context, inbound []byte) error {
	conf := &conf.InboundDetourConfig{}
	if err := json.Unmarshal(inbound, conf); err != nil {
		return fmt.Errorf("inbound unmarshal failed: %w", err)
	}

	config, err := conf.Build()
	if err != nil {
		return fmt.Errorf("inbound build failed: %w", err)
	}

	inboundConfig := command.AddInboundRequest{Inbound: config}
	_, err = c.hsClient.AddInbound(ctx, &inboundConfig)
	return handleXrayError(err)
}

func (c *Client) RemoveInbound(ctx context.Context, tag string) error {
	_, err := c.hsClient.RemoveInbound(ctx, &command.RemoveInboundRequest{
		Tag: tag,
	})
	return handleXrayError(err)
}

func (c *Client) AddUser(ctx context.Context, tag, email string, account satrapv1.Account) error {
	_, err := c.hsClient.AlterInbound(ctx, &command.AlterInboundRequest{
		Tag: tag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Email:   email,
				Account: account.ToTypedMessage(),
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

func handleXrayError(err error) error {
	s, ok := status.FromError(err)
	if !ok {
		return err
	}

	if s.Code() == codes.Unknown {
		message := s.Message()
		if strings.Contains(message, "existing tag found") {
			return errs.ErrConflict
		}
		if strings.Contains(message, "already exists") {
			return errs.ErrConflict
		}
		if strings.Contains(message, "not enough information for making a decision") {
			return errs.ErrNotFound
		}
		if strings.Contains(message, "handler not found") {
			return errs.ErrNotFound
		}
		if strings.Contains(message, "not found") {
			return errs.ErrNotFound
		}
	}

	return err
}
