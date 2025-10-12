package client

import (
	"context"
	"fmt"
	"time"

	xrayconfigv1 "github.com/vayzur/apadana/pkg/satrap/xray/config/v1"
	"github.com/xtls/xray-core/app/proxyman/command"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn                  *grpc.ClientConn
	hsClient              command.HandlerServiceClient
	runtimeRequestTimeout time.Duration
}

func New(cfg *xrayconfigv1.XrayConfig) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:                  conn,
		hsClient:              command.NewHandlerServiceClient(conn),
		runtimeRequestTimeout: cfg.RuntimeRequestTimeout,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, c.runtimeRequestTimeout)
}
