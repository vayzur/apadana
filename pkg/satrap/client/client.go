package satrap

import (
	"time"

	"github.com/vayzur/apadana/pkg/httputil"
)

type Client struct {
	httpClient *httputil.Client
}

func New(timeout time.Duration) *Client {
	httpClient := httputil.New(timeout)
	return &Client{httpClient: httpClient}
}
