package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	zlog "github.com/rs/zerolog/log"
	metav1 "github.com/vayzur/apadana/pkg/apis/meta/v1"
	satrapv1 "github.com/vayzur/apadana/pkg/apis/satrap/v1"
	"github.com/vayzur/apadana/pkg/errs"
)

func (c *Client) CreateInbound(nodeName string, inbound *satrapv1.Inbound) error {
	if nodeName == "" {
		return errs.ErrInvalidNode
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds", c.address, nodeName)
	status, resp, err := c.httpClient.Do(http.MethodPost, url, c.token, inbound)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "create").Str("nodeName", nodeName).Msg("failed")
		return err
	}

	if status == http.StatusCreated {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "create").Str("nodeName", nodeName).Int("status", status).Str("resp", string(resp)).Msg("failed")

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
				"nodeName": nodeName,
				"status":   strconv.Itoa(status),
				"resp":     string(resp),
			},
			nil,
		)
	}
}

func (c *Client) DeleteInbound(nodeName, tag string) error {
	if nodeName == "" {
		return errs.ErrInvalidNode
	}
	if tag == "" {
		return errs.ErrInvalidInbound
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s", c.address, nodeName, tag)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "delete").Str("nodeName", nodeName).Str("tag", tag).Msg("failed")
		return err
	}

	if status == http.StatusNoContent {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "delete").Str("nodeName", nodeName).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return errs.ErrInboundNotFound
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"delete inbound failed",
			map[string]string{
				"nodeName": nodeName,
				"tag":      tag,
				"status":   strconv.Itoa(status),
				"resp":     string(resp),
			},
			nil,
		)
	}
}

func (c *Client) GetInbound(nodeName, tag string) (*satrapv1.Inbound, error) {
	if nodeName == "" {
		return nil, errs.ErrInvalidNode
	}
	if tag == "" {
		return nil, errs.ErrInvalidInbound
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s", c.address, nodeName, tag)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "get").Str("nodeName", nodeName).Str("tag", tag).Msg("failed")
		return nil, err
	}

	if status == http.StatusOK {
		inbound := &satrapv1.Inbound{}
		if err := json.Unmarshal(resp, inbound); err != nil {
			zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbound").Str("action", "get").Str("nodeName", nodeName).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
			return nil, errs.New(
				errs.KindInternal,
				errs.ReasonUnmarshalFailed,
				"unmarshal inbound failed",
				map[string]string{
					"nodeName": nodeName,
					"tag":      tag,
					"status":   strconv.Itoa(status),
					"resp":     string(resp),
				},
				nil,
			)
		}
		return inbound, nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "get").Str("nodeName", nodeName).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return nil, errs.ErrInboundNotFound
	default:
		return nil, errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"get inbound failed",
			map[string]string{
				"nodeName": nodeName,
				"tag":      tag,
				"status":   strconv.Itoa(status),
				"resp":     string(resp),
			},
			nil,
		)
	}
}

func (c *Client) CountInbounds(nodeName string) (*satrapv1.Count, error) {
	if nodeName == "" {
		return nil, errs.ErrInvalidNode
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/count", c.address, nodeName)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "count").Str("nodeName", nodeName).Msg("failed")
		return nil, err
	}

	if status == http.StatusOK {
		count := &satrapv1.Count{}
		if err := json.Unmarshal(resp, count); err != nil {
			zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbound").Str("action", "count").Str("nodeName", nodeName).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
			return nil, errs.New(
				errs.KindInternal,
				errs.ReasonUnmarshalFailed,
				"inbound unmarshal failed",
				map[string]string{
					"nodeName": nodeName,
					"status":   strconv.Itoa(status),
					"resp":     string(resp),
				},
				nil,
			)
		}
		return count, nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "count").Str("nodeName", nodeName).Int("status", status).Str("resp", string(resp)).Msg("failed")

	return nil, errs.New(
		errs.KindInternal,
		errs.ReasonUnknown,
		"count inbounds failed",
		map[string]string{
			"nodeName": nodeName,
			"status":   strconv.Itoa(status),
			"resp":     string(resp),
		},
		nil,
	)
}

func (c *Client) GetInbounds(nodeName string) ([]*satrapv1.Inbound, error) {
	if nodeName == "" {
		return nil, errs.ErrInvalidNode
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds", c.address, nodeName)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbounds").Str("action", "list").Str("nodeName", nodeName).Msg("failed")
		return nil, err
	}

	if status == http.StatusOK {
		inbounds := []*satrapv1.Inbound{}
		if err := json.Unmarshal(resp, &inbounds); err != nil {
			zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbounds").Str("action", "list").Str("nodeName", nodeName).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
			return nil, errs.New(
				errs.KindInternal,
				errs.ReasonUnknown,
				"inbounds unmarshal failed",
				map[string]string{
					"nodeName": nodeName,
					"status":   strconv.Itoa(status),
					"resp":     string(resp),
				},
				nil,
			)
		}
		return inbounds, nil
	}

	zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inbounds").Str("action", "list").Str("nodeName", nodeName).Int("status", status).Str("resp", string(resp)).Msg("failed")

	return nil, errs.New(
		errs.KindInternal,
		errs.ReasonUnknown,
		"get inbounds failed",
		map[string]string{
			"nodeName": nodeName,
			"status":   strconv.Itoa(status),
			"resp":     string(resp),
		},
		nil,
	)
}

func (c *Client) UpdateInboundMetadata(nodeName, tag string, newMetadata *metav1.ObjectMeta) error {
	if nodeName == "" {
		return errs.ErrInvalidNode
	}
	if tag == "" {
		return errs.ErrInvalidInbound
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/metadata", c.address, nodeName, tag)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, newMetadata)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "update").Str("nodeName", nodeName).Str("tag", tag).Msg("failed")
		return err
	}

	if status == http.StatusOK {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "update").Str("nodeName", nodeName).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return errs.ErrInboundNotFound
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"update inbound metadata failed",
			map[string]string{
				"nodeName": nodeName,
				"tag":      tag,
				"status":   strconv.Itoa(status),
				"resp":     string(resp),
			},
			nil,
		)
	}
}

func (c *Client) UpdateInboundSpec(nodeName, tag string, newSpec *satrapv1.InboundSpec) error {
	if nodeName == "" {
		return errs.ErrInvalidNode
	}
	if tag == "" {
		return errs.ErrInvalidInbound
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/spec", c.address, nodeName, tag)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, newSpec)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inbound").Str("action", "update").Str("nodeName", nodeName).Str("tag", tag).Msg("failed")
		return err
	}

	if status == http.StatusOK {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inbound").Str("action", "update").Str("nodeName", nodeName).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return errs.ErrInboundNotFound
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"update inbound spec failed",
			map[string]string{
				"nodeName": nodeName,
				"tag":      tag,
				"status":   strconv.Itoa(status),
				"resp":     string(resp),
			},
			nil,
		)
	}
}

func (c *Client) GetInboundUsers(nodeName, tag string) ([]*satrapv1.InboundUser, error) {
	if nodeName == "" {
		return nil, errs.ErrInvalidNode
	}
	if tag == "" {
		return nil, errs.ErrInvalidInbound
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users", c.address, nodeName, tag)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUsers").Str("action", "list").Str("nodeName", nodeName).Str("tag", tag).Msg("failed")
		return nil, err
	}

	if status == http.StatusOK {
		inboundUsers := []*satrapv1.InboundUser{}
		if err := json.Unmarshal(resp, &inboundUsers); err != nil {
			zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "list").Str("nodeName", nodeName).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")
			return nil, errs.New(
				errs.KindInternal,
				errs.ReasonUnknown,
				"inbound users unmarshal failed",
				map[string]string{
					"nodeName": nodeName,
					"tag":      tag,
					"status":   strconv.Itoa(status),
					"resp":     string(resp),
				},
				nil,
			)
		}
		return inboundUsers, nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "list").Str("nodeName", nodeName).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")

	return nil, errs.New(
		errs.KindInternal,
		errs.ReasonUnknown,
		"get inbound users failed",
		map[string]string{
			"nodeName": nodeName,
			"tag":      tag,
			"status":   strconv.Itoa(status),
			"resp":     string(resp),
		},
		nil,
	)
}

func (c *Client) CreateInboundUser(nodeName, tag string, user *satrapv1.InboundUser) error {
	if nodeName == "" {
		return errs.ErrInvalidNode
	}
	if tag == "" {
		return errs.ErrInvalidInbound
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users", c.address, nodeName, tag)
	status, resp, err := c.httpClient.Do(http.MethodPost, url, c.token, user)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUsers").Str("action", "create").Str("nodeName", nodeName).Str("tag", tag).Msg("failed")
		return err
	}

	if status == http.StatusCreated {
		return nil
	}

	zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "create").Str("nodeName", nodeName).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")

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
				"nodeName": nodeName,
				"tag":      tag,
				"status":   strconv.Itoa(status),
				"resp":     string(resp),
			},
			nil,
		)
	}
}

func (c *Client) DeleteInboundUser(nodeName, tag, email string) error {
	if nodeName == "" {
		return errs.ErrInvalidNode
	}
	if tag == "" {
		return errs.ErrInvalidInbound
	}
	if email == "" {
		return errs.ErrInvalidUser
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users/%s", c.address, nodeName, tag, email)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUsers").Str("action", "delete").Str("nodeName", nodeName).Str("tag", tag).Str("email", email).Msg("failed")
		return err
	}

	if status == http.StatusNoContent {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inboundUsers").Str("action", "delete").Str("nodeName", nodeName).Str("tag", tag).Str("email", email).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return errs.ErrUserNotFound
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"delete inbound user failed",
			map[string]string{
				"nodeName": nodeName,
				"tag":      tag,
				"email":    email,
				"status":   strconv.Itoa(status),
				"resp":     string(resp),
			},
			nil,
		)
	}

}

func (c *Client) UpdateInboundUserMetadata(nodeName, tag, email string, newMetadata *metav1.ObjectMeta) error {
	if nodeName == "" {
		return errs.ErrInvalidNode
	}
	if tag == "" {
		return errs.ErrInvalidInbound
	}
	if email == "" {
		return errs.ErrInvalidUser
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users/%s/spec", c.address, nodeName, tag, email)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, newMetadata)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUser").Str("action", "update").Str("nodeName", nodeName).Str("tag", tag).Str("email", email).Msg("failed")
		return err
	}

	if status == http.StatusOK {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inboundUser").Str("action", "update").Str("nodeName", nodeName).Str("tag", tag).Str("email", email).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return errs.ErrUserNotFound
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"update inbound user metadata failed",
			map[string]string{
				"nodeName": nodeName,
				"tag":      tag,
				"email":    email,
				"status":   strconv.Itoa(status),
				"resp":     string(resp),
			},
			nil,
		)
	}
}

func (c *Client) UpdateInboundUserSpec(nodeName, tag, email string, newSpec *satrapv1.InboundUserSpec) error {
	if nodeName == "" {
		return errs.ErrInvalidNode
	}
	if tag == "" {
		return errs.ErrInvalidInbound
	}
	if email == "" {
		return errs.ErrInvalidUser
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users/%s/spec", c.address, nodeName, tag, email)
	status, resp, err := c.httpClient.Do(http.MethodPatch, url, c.token, newSpec)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUser").Str("action", "update").Str("nodeName", nodeName).Str("tag", tag).Str("email", email).Msg("failed")
		return err
	}

	if status == http.StatusOK {
		return nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inboundUser").Str("action", "update").Str("nodeName", nodeName).Str("tag", tag).Str("email", email).Int("status", status).Str("resp", string(resp)).Msg("failed")

	switch status {
	case http.StatusNotFound:
		return errs.ErrUserNotFound
	default:
		return errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"update inbound user spec failed",
			map[string]string{
				"nodeName": nodeName,
				"tag":      tag,
				"email":    email,
				"status":   strconv.Itoa(status),
				"resp":     string(resp),
			},
			nil,
		)
	}
}

func (c *Client) CountInboundUsers(nodeName, tag string) (*satrapv1.Count, error) {
	if nodeName == "" {
		return nil, errs.ErrInvalidNode
	}
	if tag == "" {
		return nil, errs.ErrInvalidInbound
	}
	url := fmt.Sprintf("%s/api/v1/nodes/%s/inbounds/%s/users/count", c.address, nodeName, tag)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, c.token, nil)
	if err != nil {
		zlog.Error().Err(err).Str("component", "client").Str("resource", "inboundUser").Str("action", "count").Str("nodeName", nodeName).Str("tag", tag).Msg("failed")
		return nil, err
	}

	if status == http.StatusOK {
		count := &satrapv1.Count{}
		if err := json.Unmarshal(resp, count); err != nil {
			zlog.Error().Err(err).Str("component", "apadana").Str("resource", "inboundUser").Str("action", "count").Str("nodeName", nodeName).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("unmarshal failed")

		}
		return count, nil
	}

	zlog.Error().Str("component", "apadana").Str("resource", "inboundUser").Str("action", "count").Str("nodeName", nodeName).Str("tag", tag).Int("status", status).Str("resp", string(resp)).Msg("failed")
	switch status {
	case http.StatusNotFound:
		return nil, errs.ErrUserNotFound
	default:
		return nil, errs.New(
			errs.KindInternal,
			errs.ReasonUnknown,
			"count inbound users failed",
			map[string]string{
				"nodeName": nodeName,
				"tag":      tag,
				"status":   strconv.Itoa(status),
				"resp":     string(resp),
			},
			nil,
		)
	}
}
