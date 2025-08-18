package satrap

import "github.com/vayzur/apadana/pkg/httputil"

type Client struct {
	httpClient *httputil.Client
}

func New(httpClient *httputil.Client) *Client {
	return &Client{httpClient: httpClient}
}
