package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	zlog "github.com/rs/zerolog/log"
	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
	"github.com/vayzur/apadana/pkg/errs"
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
		zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "create").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.ErrConflict
	}
	if status == http.StatusTooManyRequests {
		zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "create").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.ErrCapacityExceeded
	}
	if status != http.StatusCreated {
		zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "create").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
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
		zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "delete").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.ErrNotFound
	}
	if status != http.StatusNoContent {
		zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "delete").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("tag", tag).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
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
		zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "get").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return nil, errs.ErrNotFound
	}
	if status != http.StatusOK {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbound").Str("action", "get").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return nil, errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("tag", tag).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}

	inbound := &satrapv1.Inbound{}
	if err := json.Unmarshal(resp, inbound); err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbound").Str("action", "get").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
		return nil, errs.New("unmarshal", "unmarshal failed").WithField("nodeID", nodeID).WithField("tag", tag).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}

	return inbound, nil
}

func (c *Client) CountInbounds(nodeID string) (*satrapv1.Count, error) {
	if nodeID == "" {
		return nil, fmt.Errorf("nodeID cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/count", c.address, nodeID)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "count").Str("nodeID", nodeID).Msg("failed")
		return nil, err
	}
	if status == http.StatusNotFound {
		zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "count").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return nil, errs.ErrNotFound
	}
	if status != http.StatusOK {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbound").Str("action", "count").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return nil, errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}

	count := &satrapv1.Count{}
	if err := json.Unmarshal(resp, count); err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbound").Str("action", "count").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
		return nil, errs.New("unmarshal", "unmarshal failed").WithField("nodeID", nodeID).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}

	return count, nil
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
		zlog.Error().Str("component", "apadana").Str("resource", "inbounds").Str("action", "list").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return nil, errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	inbounds := []*satrapv1.Inbound{}
	if err := json.Unmarshal(resp, &inbounds); err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbounds").Str("action", "list").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
		return nil, errs.New("unmarshal", "unmarshal failed").WithField("nodeID", nodeID).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return inbounds, nil
}

func (c *Client) RenewInbound(nodeID, tag string, renew *satrapv1.Renew) error {
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}
	if tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/renew", c.address, nodeID, tag)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, renew)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Msg("failed")
		return err
	}
	if status == http.StatusNotFound {
		zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.ErrNotFound
	}
	if status != http.StatusOK {
		zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("tag", tag).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return nil
}

func (c *Client) UpdateInboundMetadata(nodeID, tag string, metadata *satrapv1.Metadata) error {
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}
	if tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/renew", c.address, nodeID, tag)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, metadata)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Msg("failed")
		return err
	}
	if status == http.StatusNotFound {
		zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.ErrNotFound
	}
	if status != http.StatusOK {
		zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("tag", tag).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return nil
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
		zlog.Error().Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "list").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return nil, errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("tag", tag).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	inboundUsers := []*satrapv1.InboundUser{}
	if err := json.Unmarshal(resp, &inboundUsers); err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "list").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
		return nil, errs.New("unmarshal", "unmarshal failed").WithField("nodeID", nodeID).WithField("tag", tag).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
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
		return errs.ErrConflict
	}
	if status != http.StatusCreated {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "create").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("tag", tag).WithField("email", user.Email).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
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
		zlog.Error().Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "delete").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.ErrNotFound
	}
	if status != http.StatusNoContent {
		zlog.Error().Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "delete").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("tag", tag).WithField("email", email).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return nil
}

func (c *Client) RenewInboundUser(nodeID, tag, email string, renew *satrapv1.Renew) error {
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}
	if tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users/%s/renew", c.address, nodeID, tag, email)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, renew)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUser").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Msg("failed")
		return err
	}
	if status == http.StatusNotFound {
		zlog.Error().Str("component", "apadana").Str("resource", "inboundUser").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.ErrNotFound
	}
	if status != http.StatusOK {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUser").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("tag", tag).WithField("email", email).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return nil
}

func (c *Client) UpdateInboundUserMetadata(nodeID, tag, email string, metadata *satrapv1.Metadata) error {
	if nodeID == "" {
		return fmt.Errorf("nodeID cannot be empty")
	}
	if tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users/%s/renew", c.address, nodeID, tag, email)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, metadata)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUser").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Msg("failed")
		return err
	}
	if status == http.StatusNotFound {
		zlog.Error().Str("component", "apadana").Str("resource", "inboundUser").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.ErrNotFound
	}
	if status != http.StatusOK {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUser").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("tag", tag).WithField("email", email).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return nil
}

func (c *Client) CountInboundUsers(nodeID, tag string) (*satrapv1.Count, error) {
	if nodeID == "" {
		return nil, fmt.Errorf("nodeID cannot be empty")
	}
	if tag == "" {
		return nil, fmt.Errorf("tag cannot be empty")
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users/count", c.address, nodeID, tag)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUser").Str("action", "count").Str("nodeID", nodeID).Str("tag", tag).Msg("failed")
		return nil, err
	}
	if status == http.StatusNotFound {
		zlog.Error().Str("component", "apadana").Str("resource", "inboundUser").Str("action", "count").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return nil, errs.ErrNotFound
	}
	if status != http.StatusOK {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUser").Str("action", "count").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
		return nil, errs.New("unexpected", "unexpected response").WithField("nodeID", nodeID).WithField("tag", tag).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}

	count := &satrapv1.Count{}
	if err := json.Unmarshal(resp, count); err != nil {
		zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUser").Str("action", "count").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
		return nil, errs.New("unmarshal", "unmarshal failed").WithField("nodeID", nodeID).WithField("tag", tag).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}

	return count, nil
}
