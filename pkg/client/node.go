package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/vayzur/apadana/pkg/errs"

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
		zlog.Error().Err(err).Str("component", "client").Str("resource", "node").Str("action", "update").Str("nodeID", nodeID).Msg("failed")
		return err
	}
	if status != http.StatusOK {
		zlog.Error().Str("component", "apadana").Str("resource", "node").Str("action", "update").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return nil
}

func (c *Client) UpdateNodeMetadata(nodeID string, nodeMetadata *corev1.NodeMetadata) error {
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/metadata", c.address, nodeID)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, nodeMetadata)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "node").Str("action", "update").Str("nodeID", nodeID).Msg("failed")
		return err
	}
	if status != http.StatusOK {
		zlog.Error().Str("component", "apadana").Str("resource", "node").Str("action", "update").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return nil
}

func (c *Client) UpdateNodeSpec(nodeID string, nodeSpec *corev1.NodeSpec) error {
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/spec", c.address, nodeID)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, nodeSpec)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "node").Str("action", "update").Str("nodeID", nodeID).Msg("failed")
		return err
	}
	if status != http.StatusOK {
		zlog.Error().Str("component", "apadana").Str("resource", "node").Str("action", "update").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return nil
}

func (c *Client) GetNode(nodeID string) (*corev1.Node, error) {
	if nodeID == "" {
		return nil, fmt.Errorf("nodeID cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s", c.address, nodeID)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "node").Str("action", "get").Str("nodeID", nodeID).Msg("failed")
		return nil, err
	}
	if status != http.StatusOK {
		zlog.Error().Str("component", "apadana").Str("resource", "node").Str("action", "get").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return nil, errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	node := &corev1.Node{}
	if err := json.Unmarshal(resp, node); err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "node").Str("action", "get").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
		return nil, errs.New("unmarshal", "unmarshal failed").WithField("nodeID", nodeID).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return node, nil
}

func (c *Client) GetNodes() ([]*corev1.Node, error) {
	url := fmt.Sprintf("%s/api/v1/nodes", c.address)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "nodes").Str("action", "list").Msg("failed")
		return nil, err
	}
	if status != http.StatusOK {
		zlog.Error().Str("component", "apadana").Str("resource", "nodes").Str("action", "list").Int("status", status).Str("resp", string(resp)).Msg("failed")
		return nil, errs.New("unexpected", "unexpected response").WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	nodes := []*corev1.Node{}
	if err := json.Unmarshal(resp, &nodes); err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "nodes").Str("action", "list").Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
		return nil, errs.New("unmarshal", "unmarshal failed").WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return nodes, nil
}

func (c *Client) GetActiveNodes() ([]*corev1.Node, error) {
	url := fmt.Sprintf("%s/api/v1/nodes/active", c.address)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "nodes").Str("action", "list").Msg("failed")
		return nil, err
	}
	if status != http.StatusOK {
		zlog.Error().Str("component", "apadana").Str("resource", "nodes").Str("action", "list").Int("status", status).Str("resp", string(resp)).Msg("failed")
		return nil, errs.New("unexpected", "unexpected response").WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	nodes := []*corev1.Node{}
	if err := json.Unmarshal(resp, &nodes); err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "nodes").Str("action", "list").Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
		return nil, errs.New("unmarshal", "unmarshal failed").WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return nodes, nil
}

func (c *Client) CreateNode(node *corev1.Node) (*corev1.Node, error) {
	url := fmt.Sprintf("%s/api/v1/nodes", c.address)
	status, resp, err := c.httpClient.Do(http.MethodPost, url, c.token, node)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "node").Str("action", "create").Msg("failed")
		return nil, err
	}
	if status != http.StatusCreated {
		zlog.Error().Str("component", "apadana").Str("resource", "node").Str("action", "create").Int("status", status).Str("resp", string(resp)).Msg("failed")
		return nil, errs.New("unexpected", "unexpected response").WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	n := &corev1.Node{}
	if err := json.Unmarshal(resp, n); err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "node").Str("action", "create").Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
		return nil, errs.New("unmarshal", "unmarshal failed").WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return n, nil
}

func (c *Client) DeleteNode(nodeID string) error {
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s", c.address, nodeID)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "node").Str("action", "delete").Str("nodeID", nodeID).Msg("failed")
		return err
	}
	if status == http.StatusNotFound {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "node").Str("action", "delete").Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.ErrNotFound
	}
	if status != http.StatusNoContent {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "node").Str("action", "delete").Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return nil
}
