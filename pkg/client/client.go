package client

import (
	"time"

	"github.com/vayzur/apadana/pkg/httputil"
)

type Client struct {
	httpClient *httputil.Client
	address    string
	token      string
}

func New(address, token string, timeout time.Duration) *Client {
	httpClient := httputil.New(timeout)
	return &Client{
		httpClient: httpClient,
		address:    address,
		token:      token,
	}
}
