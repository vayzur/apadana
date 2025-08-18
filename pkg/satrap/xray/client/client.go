package client

import (
	"fmt"

	"github.com/xtls/xray-core/app/proxyman/command"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn     *grpc.ClientConn
	hsClient command.HandlerServiceClient
}

func New(endpoint string) (*Client, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
