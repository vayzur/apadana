package satrap

import (
	"fmt"
	"net/http"
	"strconv"

	corev1 "github.com/vayzur/apadana/pkg/apis/core/v1"
	satrapv1 "github.com/vayzur/apadana/pkg/apis/satrap/v1"
	"github.com/vayzur/apadana/pkg/errs"
	"github.com/xtls/xray-core/infra/conf"
)

func (c *Client) AddInbound(node *corev1.Node, inboundConfig *conf.InboundDetourConfig) error {
	url := node.URL("/api/v1/inbounds")
	status, resp, err := c.httpClient.Do(http.MethodPost, url, node.Spec.Token, inboundConfig)
	if err != nil {
		return err
	}
	if status == http.StatusCreated {
		return nil
	}
	if status == http.StatusConflict {
		return errs.ErrInboundConflict
	}
	return errs.New(
		errs.KindInternal,
		errs.ReasonUnknown,
		"add inbound failed",
		map[string]string{
			"nodeName": node.Metadata.Name,
			"status":   strconv.Itoa(status),
			"resp":     string(resp),
		},
		nil,
	)
}

func (c *Client) RemoveInbound(node *corev1.Node, tag string) error {
	if tag == "" {
		return errs.ErrInvalidInbound
	}
	path := fmt.Sprintf("/api/v1/inbounds/%s", tag)
	url := node.URL(path)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, node.Spec.Token, nil)
	if err != nil {
		return err
	}
	if status == http.StatusNoContent {
		return nil
	}
	if status == http.StatusNotFound {
		return errs.ErrInboundNotFound
	}
	return errs.New(
		errs.KindInternal,
		errs.ReasonUnknown,
		"remove inbound failed",
		map[string]string{
			"nodeName": node.Metadata.Name,
			"tag":      tag,
			"status":   strconv.Itoa(status),
			"resp":     string(resp),
		},
		nil,
	)
}

func (c *Client) AddUser(node *corev1.Node, tag string, user *satrapv1.InboundUser) error {
	if tag == "" {
		return errs.ErrInvalidInbound
	}
	path := fmt.Sprintf("/api/v1/inbounds/%s/users", tag)
	url := node.URL(path)
	status, resp, err := c.httpClient.Do(http.MethodPost, url, node.Spec.Token, user)
	if err != nil {
		return err
	}
	if status == http.StatusCreated {
		return nil
	}
	if status == http.StatusConflict {
		return errs.ErrUserConflict
	}
	return errs.New(
		errs.KindInternal,
		errs.ReasonUnknown,
		"add user failed",
		map[string]string{
			"nodeName": node.Metadata.Name,
			"tag":      tag,
			"status":   strconv.Itoa(status),
			"resp":     string(resp),
		},
		nil,
	)
}

func (c *Client) RemoveUser(node *corev1.Node, tag, email string) error {
	if tag == "" {
		return errs.ErrInvalidInbound
	}
	if email == "" {
		return errs.ErrInvalidUser
	}
	path := fmt.Sprintf("/api/v1/inbounds/%s/users/%s", tag, email)
	url := node.URL(path)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, node.Spec.Token, nil)
	if err != nil {
		return err
	}
	if status == http.StatusNoContent {
		return nil
	}
	if status == http.StatusNotFound {
		return errs.ErrUserNotFound
	}
	return errs.New(
		errs.KindInternal,
		errs.ReasonUnknown,
		"remove user failed",
		map[string]string{
			"nodeName": node.Metadata.Name,
			"tag":      tag,
			"email":    email,
			"status":   strconv.Itoa(status),
			"resp":     string(resp),
		},
		nil,
	)
}
