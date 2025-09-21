package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/vayzur/apadana/pkg/errs"

	zlog "github.com/rs/zerolog/log"
	corev1 "github.com/vayzur/apadana/pkg/apis/core/v1"
	metav1 "github.com/vayzur/apadana/pkg/apis/meta/v1"
)

func (c *Client) UpdateNodeStatus(nodeName string, nodeStatus *corev1.NodeStatus) error {
	if nodeName == "" {
		return errs.ErrInvalidNode
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/status", c.address, nodeName)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, nodeStatus)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "node").Str("action", "update").Str("nodeName", nodeName).Msg("failed")
		return err
	}

	if status == http.StatusOK {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "node").Str("action", "update").Str("nodeName", nodeName).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return errs.ErrNodeNotFound
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"update node status failed",
			map[string]string{
				"nodeName": nodeName,
				"status":   strconv.Itoa(status),
				"resp":     string(resp),
			},
			nil,
		)
	}

}

func (c *Client) UpdateNodeMetadata(nodeName string, nodeMetadata *metav1.ObjectMeta) error {
	if nodeName == "" {
		return errs.ErrInvalidNode
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/metadata", c.address, nodeName)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, nodeMetadata)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "node").Str("action", "update").Str("nodeName", nodeName).Msg("failed")
		return err
	}

	if status == http.StatusOK {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "node").Str("action", "update").Str("nodeName", nodeName).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return errs.ErrNodeNotFound
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"update node metadata failed",
			map[string]string{
				"nodeName": nodeName,
				"status":   strconv.Itoa(status),
				"resp":     string(resp),
			},
			nil,
		)
	}
}

func (c *Client) UpdateNodeSpec(nodeName string, nodeSpec *corev1.NodeSpec) error {
	if nodeName == "" {
		return errs.ErrInvalidNode
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/spec", c.address, nodeName)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, nodeSpec)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "node").Str("action", "update").Str("nodeName", nodeName).Msg("failed")
		return err
	}

	if status == http.StatusOK {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "node").Str("action", "update").Str("nodeName", nodeName).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return errs.ErrNodeNotFound
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"update node spec failed",
			map[string]string{
				"nodeName": nodeName,
				"status":   strconv.Itoa(status),
				"resp":     string(resp),
			},
			nil,
		)
	}
}

func (c *Client) GetNode(nodeName string) (*corev1.Node, error) {
	if nodeName == "" {
		return nil, errs.ErrInvalidNode
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s", c.address, nodeName)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "node").Str("action", "get").Str("nodeName", nodeName).Msg("failed")
		return nil, err
	}

	if status == http.StatusOK {
		node := &corev1.Node{}
		if err := json.Unmarshal(resp, node); err != nil {
			zlog.Error().Err(err).Str("component", "apadana").Str("resource", "node").Str("action", "get").Str("nodeName", nodeName).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
			return nil, errs.New(
				errs.KindInternal,
				errs.ReasonUnmarshalFailed,
				"node unmarshal failed",
				map[string]string{
					"nodeName": nodeName,
					"status":   strconv.Itoa(status),
					"resp":     string(resp),
				},
				nil,
			)
		}
		return node, nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "node").Str("action", "get").Str("nodeName", nodeName).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return nil, errs.ErrNodeNotFound
	default:
		return nil, errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"get node failed",
			map[string]string{
				"nodeName": nodeName,
				"status":   strconv.Itoa(status),
				"resp":     string(resp),
			},
			nil,
		)
	}
}

func (c *Client) GetNodes() ([]*corev1.Node, error) {
	url := fmt.Sprintf("%s/api/v1/nodes", c.address)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "nodes").Str("action", "list").Msg("failed")
		return nil, err
	}
	if status == http.StatusOK {
		nodes := []*corev1.Node{}
		if err := json.Unmarshal(resp, &nodes); err != nil {
			zlog.Error().Err(err).Str("component", "apadana").Str("resource", "nodes").Str("action", "list").Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
			return nil, errs.New(
				errs.KindInternal,
				errs.ReasonUnmarshalFailed,
				"nodes unmarshal failed",
				map[string]string{
					"status": strconv.Itoa(status),
					"resp":   string(resp),
				},
				nil,
			)
		}
		return nodes, nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "nodes").Str("action", "list").Int("status", status).Str("resp", string(resp)).Msg("failed")

	return nil, errs.New(
		errs.KindInternal,
		errs.ReasonUnknown,
		"get nodes failed",
		map[string]string{
			"status": strconv.Itoa(status),
			"resp":   string(resp),
		},
		nil,
	)
}

func (c *Client) GetActiveNodes() ([]*corev1.Node, error) {
	url := fmt.Sprintf("%s/api/v1/nodes/active", c.address)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "nodes").Str("action", "list").Msg("failed")
		return nil, err
	}
	if status == http.StatusOK {
		nodes := []*corev1.Node{}
		if err := json.Unmarshal(resp, &nodes); err != nil {
			zlog.Error().Err(err).Str("component", "apadana").Str("resource", "nodes").Str("action", "list").Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
			return nil, errs.New(
				errs.KindInternal,
				errs.ReasonUnmarshalFailed,
				"active nodes unmarshal failed",
				map[string]string{
					"status": strconv.Itoa(status),
					"resp":   string(resp),
				},
				nil,
			)
		}
		return nodes, nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "nodes").Str("action", "list").Int("status", status).Str("resp", string(resp)).Msg("failed")

	return nil, errs.New(
		errs.KindInternal,
		errs.ReasonUnknown,
		"get active nodes failed",
		map[string]string{
			"status": strconv.Itoa(status),
			"resp":   string(resp),
		},
		nil,
	)
}

func (c *Client) CreateNode(node *corev1.Node) (*corev1.Node, error) {
	url := fmt.Sprintf("%s/api/v1/nodes", c.address)
	status, resp, err := c.httpClient.Do(http.MethodPost, url, c.token, node)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "node").Str("action", "create").Msg("failed")
		return nil, err
	}

	if status == http.StatusCreated {
		n := &corev1.Node{}
		if err := json.Unmarshal(resp, n); err != nil {
			zlog.Error().Err(err).Str("component", "apadana").Str("resource", "node").Str("action", "create").Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
			return nil, errs.New(
				errs.KindInternal,
				errs.ReasonUnmarshalFailed,
				"node unmarshal failed",
				map[string]string{
					"status": strconv.Itoa(status),
					"resp":   string(resp),
				},
				nil,
			)
		}
		return n, nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "node").Str("action", "create").Int("status", status).Str("resp", string(resp)).Msg("failed")

	return nil, errs.New(
		errs.KindInternal,
		errs.ReasonUnknown,
		"create node failed",
		map[string]string{
			"status": strconv.Itoa(status),
			"resp":   string(resp),
		},
		nil,
	)
}

func (c *Client) DeleteNode(nodeName string) error {
	if nodeName == "" {
		return errs.ErrInvalidNode
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s", c.address, nodeName)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "node").Str("action", "delete").Str("nodeName", nodeName).Msg("failed")
		return err
	}

	if status == http.StatusNoContent {
		return nil
	}

	zlog.Error().Err(err).Str("component", "apadana").Str("resource", "node").Str("action", "delete").Str("nodeName", nodeName).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return errs.ErrNodeNotFound
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"get node failed",
			map[string]string{
				"nodeName": nodeName,
				"status":   strconv.Itoa(status),
				"resp":     string(resp),
			},
			nil,
		)
	}

}
