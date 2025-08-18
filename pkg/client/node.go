package apadana

import (
	"fmt"
	"net/http"

	zlog "github.com/rs/zerolog/log"
	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
)

func (c *Client) UpdateNodeStatus(nodeID string, nodeStatus *corev1.NodeStatus) error {
	url := fmt.Sprintf("%s/api/v1/nodes/%s/status", c.address, nodeID)

	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, nodeStatus)
	if err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Msg("failed to send node update status")
		return err
	}
	if status != 200 {
		zlog.Error().Err(err).Str("component", "apadana").Str("resp", string(resp)).Int("status", status).Msg("node status update failed")
		return err
	}

	return nil
}
