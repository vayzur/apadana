package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	zlog "github.com/rs/zerolog/log"
	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
)

func (c *Client) UpdateNodeStatus(nodeID string, nodeStatus *corev1.NodeStatus) error {
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}

	url := fmt.Sprintf("%s/api/v1/nodes/%s/status", c.address, nodeID)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, nodeStatus)
	if err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Msg("failed to send node update status")
		return err
	}
	if status != http.StatusOK {
		zlog.Error().Err(err).Str("component", "apadana").Str("resp", string(resp)).Int("status", status).Msg("node status update failed")
		return err
	}

	return nil
}

func (c *Client) GetNodes() ([]*corev1.Node, error) {
	url := fmt.Sprintf("%s/api/v1/nodes", c.address)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Msg("failed to send get nodes")
		return nil, err
	}
	if status != http.StatusOK {
		zlog.Error().Err(err).Str("component", "apadana").Str("resp", string(resp)).Int("status", status).Msg("get nodes failed")
		return nil, err
	}
	var nodes []*corev1.Node
	if err := json.Unmarshal(resp, &nodes); err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Msg("unmarshal nodes failed")
		return nil, err
	}
	return nodes, nil
}

func (c *Client) GetActiveNodes() ([]*corev1.Node, error) {
	url := fmt.Sprintf("%s/api/v1/nodes/active", c.address)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Msg("failed to send get nodes")
		return nil, err
	}
	if status != http.StatusOK {
		zlog.Error().Err(err).Str("component", "apadana").Str("resp", string(resp)).Int("status", status).Msg("get nodes failed")
		return nil, err
	}
	var nodes []*corev1.Node
	if err := json.Unmarshal(resp, &nodes); err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Msg("unmarshal nodes failed")
		return nil, err
	}
	return nodes, nil
}
