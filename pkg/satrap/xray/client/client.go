package client

import (
	"fmt"

	xrayconfigv1 "github.com/vayzur/apadana/pkg/satrap/xray/config/v1"
	"github.com/xtls/xray-core/app/proxyman/command"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn     *grpc.ClientConn
	hsClient command.HandlerServiceClient
}

func New(cfg *xrayconfigv1.XrayConfig) (*Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("xray connect failed: %w", err)
	}

	return &Client{
		conn:     conn,
		hsClient: command.NewHandlerServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
