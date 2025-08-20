package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	zlog "github.com/rs/zerolog/log"
	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
)

func (c *Client) CreateInbound(nodeID string, inbound *satrapv1.Inbound) error {
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds", c.address, nodeID)
	status, resp, err := c.httpClient.Do(http.MethodPost, url, c.token, inbound)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "create").Str("nodeID", nodeID).Msg("failed")
		return err
	}
	if status == http.StatusConflict {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "node").Str("action", "create").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return err
	}
	if status != http.StatusCreated {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "node").Str("action", "create").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return err
	}
	return nil
}

func (c *Client) DeleteInbound(nodeID, tag string) error {
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}
	if tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s", c.address, nodeID, tag)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "delete").Str("nodeID", nodeID).Str("tag", tag).Msg("failed")
		return err
	}
	if status == http.StatusNotFound {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbound").Str("action", "delete").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return err
	}
	if status != http.StatusNoContent {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "node").Str("action", "delete").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return err
	}
	return nil
}

func (c *Client) GetInbound(nodeID, tag string) (*satrapv1.Inbound, error) {
	if nodeID == "" {
		return nil, fmt.Errorf("nodeID cannot be empty")
	}
	if tag == "" {
		return nil, fmt.Errorf("tag cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s", c.address, nodeID, tag)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "get").Str("nodeID", nodeID).Str("tag", tag).Msg("failed")
		return nil, err
	}
	if status == http.StatusNotFound {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbound").Str("action", "get").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return nil, err
	}
	if status != http.StatusOK {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbound").Str("action", "get").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return nil, err
	}

	inbound := &satrapv1.Inbound{}
	if err := json.Unmarshal(resp, inbound); err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbound").Str("action", "get").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
		return nil, err
	}

	return inbound, nil
}

func (c *Client) GetInbounds(nodeID string) ([]*satrapv1.Inbound, error) {
	if nodeID == "" {
		return nil, fmt.Errorf("nodeID cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds", c.address, nodeID)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbounds").Str("action", "list").Str("nodeID", nodeID).Msg("failed")
		return nil, err
	}
	if status != http.StatusOK {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbounds").Str("action", "list").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return nil, err
	}
	inbounds := []*satrapv1.Inbound{}
	if err := json.Unmarshal(resp, &inbounds); err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbounds").Str("action", "list").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
		return nil, err
	}
	return inbounds, nil
}

func (c *Client) GetInboundUsers(nodeID, tag string) ([]*satrapv1.InboundUser, error) {
	if nodeID == "" {
		return nil, fmt.Errorf("nodeID cannot be empty")
	}
	if tag == "" {
		return nil, fmt.Errorf("tag cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users", c.address, nodeID, tag)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUsers").Str("action", "list").Str("nodeID", nodeID).Str("tag", tag).Msg("failed")
		return nil, err
	}
	if status != http.StatusOK {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "list").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return nil, err
	}
	inboundUsers := []*satrapv1.InboundUser{}
	if err := json.Unmarshal(resp, &inboundUsers); err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "list").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
		return nil, err
	}
	return inboundUsers, nil
}

func (c *Client) CreateInboundUser(nodeID, tag string, user *satrapv1.InboundUser) error {
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}
	if tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users", c.address, nodeID, tag)
	status, resp, err := c.httpClient.Do(http.MethodPost, url, c.token, user)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUsers").Str("action", "create").Str("nodeID", nodeID).Str("tag", tag).Msg("failed")
		return err
	}
	if status == http.StatusConflict {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "create").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return err
	}
	if status != http.StatusCreated {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "create").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return err
	}
	return nil
}

func (c *Client) DeleteInboundUser(nodeID, tag, email string) error {
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}
	if tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users/%s", c.address, nodeID, tag, email)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUsers").Str("action", "delete").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Msg("failed")
		return err
	}
	if status == http.StatusNotFound {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "delete").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return err
	}
	if status != http.StatusNoContent {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "delete").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return err
	}
	return nil
}
