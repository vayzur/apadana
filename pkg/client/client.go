package apadana

import "github.com/vayzur/apadana/pkg/httputil"

type Client struct {
	httpClient *httputil.Client
	address    string
	token      string
}

func New(httpClient *httputil.Client, address, token string) *Client {
	return &Client{
		httpClient: httpClient,
		address:    address,
		token:      token,
	}
}
