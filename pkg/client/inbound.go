package client

import (
	"fmt"
	"net/http"

	zlog "github.com/rs/zerolog/log"
	v1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
)

func (c *Client) AddInbound(nodeID string, inbound *v1.Inbound) error {
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds", c.address, nodeID)

	status, resp, err := c.httpClient.Do(http.MethodPost, url, c.token, inbound)
	if err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Msg("failed to send add inbound")
		return fmt.Errorf("add inbound send failed: %w", err)
	}
	if status == http.StatusConflict {
		zlog.Error().Err(err).Str("component", "apadana").Int("status", status).Msg("inbound exists")
		return fmt.Errorf("inbound exists: %w", err)
	}
	if status != http.StatusCreated {
		zlog.Error().Err(err).Str("component", "apadana").Int("status", status).Msg("add inbound failed")
		return fmt.Errorf("add inbound failed with status: %d resp %s: %w", status, string(resp), err)
	}

	return nil
}
