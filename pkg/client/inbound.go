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
		return errs.ErrInvalidNodeID
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds", c.address, nodeID)
	status, resp, err := c.httpClient.Do(http.MethodPost, url, c.token, inbound)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "create").Str("nodeID", nodeID).Msg("failed")
		return err
	}

	if status == http.StatusCreated {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "create").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusConflict:
		return errs.ErrInboundConflict
	case http.StatusTooManyRequests:
		return errs.ErrNodeCapacityExceeded
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"create inbound failed",
			map[string]string{
				"nodeID": nodeID,
				"status": strconv.Itoa(status),
				"resp":   string(resp),
			},
			nil,
		)
	}
}

func (c *Client) DeleteInbound(nodeID, tag string) error {
	if nodeID == "" {
		return errs.ErrInvalidNodeID
	}
	if tag == "" {
		return errs.ErrInvalidTag
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s", c.address, nodeID, tag)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "delete").Str("nodeID", nodeID).Str("tag", tag).Msg("failed")
		return err
	}

	if status == http.StatusNoContent {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "delete").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return errs.ErrInboundNotFound
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"delete inbound failed",
			map[string]string{
				"nodeID": nodeID,
				"tag":    tag,
				"status": strconv.Itoa(status),
				"resp":   string(resp),
			},
			nil,
		)
	}
}

func (c *Client) GetInbound(nodeID, tag string) (*satrapv1.Inbound, error) {
	if nodeID == "" {
		return nil, errs.ErrInvalidNodeID
	}
	if tag == "" {
		return nil, errs.ErrInvalidTag
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s", c.address, nodeID, tag)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "get").Str("nodeID", nodeID).Str("tag", tag).Msg("failed")
		return nil, err
	}

	if status == http.StatusOK {
		inbound := &satrapv1.Inbound{}
		if err := json.Unmarshal(resp, inbound); err != nil {
			zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbound").Str("action", "get").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
			return nil, errs.New(
				errs.KindInternal,
				errs.ReasonUnmarshalFailed,
				"unmarshal inbound failed",
				map[string]string{
					"nodeID": nodeID,
					"tag":    tag,
					"status": strconv.Itoa(status),
					"resp":   string(resp),
				},
				nil,
			)
		}
		return inbound, nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "get").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return nil, errs.ErrInboundNotFound
	default:
		return nil, errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"get inbound failed",
			map[string]string{
				"nodeID": nodeID,
				"tag":    tag,
				"status": strconv.Itoa(status),
				"resp":   string(resp),
			},
			nil,
		)
	}
}

func (c *Client) CountInbounds(nodeID string) (*satrapv1.Count, error) {
	if nodeID == "" {
		return nil, errs.ErrInvalidNodeID
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/count", c.address, nodeID)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "count").Str("nodeID", nodeID).Msg("failed")
		return nil, err
	}

	if status == http.StatusOK {
		count := &satrapv1.Count{}
		if err := json.Unmarshal(resp, count); err != nil {
			zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbound").Str("action", "count").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
			return nil, errs.New(
				errs.KindInternal,
				errs.ReasonUnmarshalFailed,
				"inbound unmarshal failed",
				map[string]string{
					"nodeID": nodeID,
					"status": strconv.Itoa(status),
					"resp":   string(resp),
				},
				nil,
			)
		}
		return count, nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "count").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("failed")

	return nil, errs.New(
		errs.KindInternal,
		errs.ReasonUnknown,
		"count inbounds failed",
		map[string]string{
			"nodeID": nodeID,
			"status": strconv.Itoa(status),
			"resp":   string(resp),
		},
		nil,
	)
}

func (c *Client) GetInbounds(nodeID string) ([]*satrapv1.Inbound, error) {
	if nodeID == "" {
		return nil, errs.ErrInvalidNodeID
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds", c.address, nodeID)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbounds").Str("action", "list").Str("nodeID", nodeID).Msg("failed")
		return nil, err
	}

	if status == http.StatusOK {
		inbounds := []*satrapv1.Inbound{}
		if err := json.Unmarshal(resp, &inbounds); err != nil {
			zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbounds").Str("action", "list").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
			return nil, errs.New(
				errs.KindInternal,
				errs.ReasonUnknown,
				"inbounds unmarshal failed",
				map[string]string{
					"nodeID": nodeID,
					"status": strconv.Itoa(status),
					"resp":   string(resp),
				},
				nil,
			)
		}
		return inbounds, nil
	}

	zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbounds").Str("action", "list").Str("nodeID", nodeID).Int("status", status).Str("resp", string(resp)).Msg("failed")

	return nil, errs.New(
		errs.KindInternal,
		errs.ReasonUnknown,
		"get inbounds failed",
		map[string]string{
			"nodeID": nodeID,
			"status": strconv.Itoa(status),
			"resp":   string(resp),
		},
		nil,
	)
}

func (c *Client) RenewInbound(nodeID, tag string, renew *satrapv1.Renew) error {
	if nodeID == "" {
		return errs.ErrInvalidNodeID
	}
	if tag == "" {
		return errs.ErrInvalidTag
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/renew", c.address, nodeID, tag)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, renew)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Msg("failed")
		return err
	}

	if status == http.StatusOK {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return errs.ErrInboundNotFound
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"renew inbound failed",
			map[string]string{
				"nodeID": nodeID,
				"tag":    tag,
				"status": strconv.Itoa(status),
				"resp":   string(resp),
			},
			nil,
		)
	}
}

func (c *Client) UpdateInboundMetadata(nodeID, tag string, metadata *satrapv1.Metadata) error {
	if nodeID == "" {
		return errs.ErrInvalidNodeID
	}
	if tag == "" {
		return errs.ErrInvalidTag
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/renew", c.address, nodeID, tag)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, metadata)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Msg("failed")
		return err
	}

	if status == http.StatusOK {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return errs.ErrInboundNotFound
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"update inbound metadata failed",
			map[string]string{
				"nodeID": nodeID,
				"tag":    tag,
				"status": strconv.Itoa(status),
				"resp":   string(resp),
			},
			nil,
		)
	}
}

func (c *Client) GetInboundUsers(nodeID, tag string) ([]*satrapv1.InboundUser, error) {
	if nodeID == "" {
		return nil, errs.ErrInvalidNodeID
	}
	if tag == "" {
		return nil, errs.ErrInvalidTag
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users", c.address, nodeID, tag)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUsers").Str("action", "list").Str("nodeID", nodeID).Str("tag", tag).Msg("failed")
		return nil, err
	}

	if status == http.StatusOK {
		inboundUsers := []*satrapv1.InboundUser{}
		if err := json.Unmarshal(resp, &inboundUsers); err != nil {
			zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "list").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
			return nil, errs.New(
				errs.KindInternal,
				errs.ReasonUnknown,
				"inbound users unmarshal failed",
				map[string]string{
					"nodeID": nodeID,
					"tag":    tag,
					"status": strconv.Itoa(status),
					"resp":   string(resp),
				},
				nil,
			)
		}
		return inboundUsers, nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "list").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")

	return nil, errs.New(
		errs.KindInternal,
		errs.ReasonUnknown,
		"get inbound users failed",
		map[string]string{
			"nodeID": nodeID,
			"tag":    tag,
			"status": strconv.Itoa(status),
			"resp":   string(resp),
		},
		nil,
	)
}

func (c *Client) CreateInboundUser(nodeID, tag string, user *satrapv1.InboundUser) error {
	if nodeID == "" {
		return errs.ErrInvalidNodeID
	}
	if tag == "" {
		return errs.ErrInvalidTag
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users", c.address, nodeID, tag)
	status, resp, err := c.httpClient.Do(http.MethodPost, url, c.token, user)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUsers").Str("action", "create").Str("nodeID", nodeID).Str("tag", tag).Msg("failed")
		return err
	}

	if status == http.StatusCreated {
		return nil
	}

	zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "create").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusConflict:
		return errs.ErrUserConflict
	case http.StatusTooManyRequests:
		return errs.ErrInboundCapacityExceeded
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"create inbound user failed",
			map[string]string{
				"nodeID": nodeID,
				"tag":    tag,
				"status": strconv.Itoa(status),
				"resp":   string(resp),
			},
			nil,
		)
	}
}

func (c *Client) DeleteInboundUser(nodeID, tag, email string) error {
	if nodeID == "" {
		return errs.ErrInvalidNodeID
	}
	if tag == "" {
		return errs.ErrInvalidTag
	}
	if email == "" {
		return errs.ErrInvalidEmail
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users/%s", c.address, nodeID, tag, email)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUsers").Str("action", "delete").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Msg("failed")
		return err
	}

	if status == http.StatusNoContent {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "delete").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return errs.ErrUserNotFound
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"delete inbound user failed",
			map[string]string{
				"nodeID": nodeID,
				"tag":    tag,
				"email":  email,
				"status": strconv.Itoa(status),
				"resp":   string(resp),
			},
			nil,
		)
	}

}

func (c *Client) RenewInboundUser(nodeID, tag, email string, renew *satrapv1.Renew) error {
	if nodeID == "" {
		return errs.ErrInvalidNodeID
	}
	if tag == "" {
		return errs.ErrInvalidTag
	}
	if email == "" {
		return errs.ErrInvalidEmail
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users/%s/renew", c.address, nodeID, tag, email)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, renew)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUser").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Msg("failed")
		return err
	}

	if status == http.StatusOK {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inboundUser").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return errs.ErrUserNotFound
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"renew inbound user failed",
			map[string]string{
				"nodeID": nodeID,
				"tag":    tag,
				"email":  email,
				"status": strconv.Itoa(status),
				"resp":   string(resp),
			},
			nil,
		)
	}

}

func (c *Client) UpdateInboundUserMetadata(nodeID, tag, email string, metadata *satrapv1.Metadata) error {
	if nodeID == "" {
		return errs.ErrInvalidNodeID
	}
	if tag == "" {
		return errs.ErrInvalidTag
	}
	if email == "" {
		return errs.ErrInvalidEmail
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users/%s/renew", c.address, nodeID, tag, email)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, metadata)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUser").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Msg("failed")
		return err
	}

	if status == http.StatusOK {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inboundUser").Str("action", "update").Str("nodeID", nodeID).Str("tag", tag).Str("email", email).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return errs.ErrUserNotFound
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"update inbound user metadata failed",
			map[string]string{
				"nodeID": nodeID,
				"tag":    tag,
				"email":  email,
				"status": strconv.Itoa(status),
				"resp":   string(resp),
			},
			nil,
		)
	}
}

func (c *Client) CountInboundUsers(nodeID, tag string) (*satrapv1.Count, error) {
	if nodeID == "" {
		return nil, errs.ErrInvalidNodeID
	}
	if tag == "" {
		return nil, errs.ErrInvalidTag
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users/count", c.address, nodeID, tag)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUser").Str("action", "count").Str("nodeID", nodeID).Str("tag", tag).Msg("failed")
		return nil, err
	}

	if status == http.StatusOK {
		count := &satrapv1.Count{}
		if err := json.Unmarshal(resp, count); err != nil {
			zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUser").Str("action", "count").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")

		}
		return count, nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inboundUser").Str("action", "count").Str("nodeID", nodeID).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
	switch status {
	case http.StatusNotFound:
		return nil, errs.ErrUserNotFound
	default:
		return nil, errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"count inbound users failed",
			map[string]string{
				"nodeID": nodeID,
				"tag":    tag,
				"status": strconv.Itoa(status),
				"resp":   string(resp),
			},
			nil,
		)
	}
}
