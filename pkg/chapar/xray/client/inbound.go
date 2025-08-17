package xray

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vayzur/apadana/pkg/errs"
	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/infra/conf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) AddInbound(ctx context.Context, inbound []byte) error {
	conf := new(conf.InboundDetourConfig)
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
		if strings.Contains(message, "not enough information for making a decision") {
			return errs.ErrNotFound
		}
	}

	return err
}
