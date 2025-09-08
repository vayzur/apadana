package satrap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
	"github.com/vayzur/apadana/pkg/errs"
)

func (c *Client) CountInbounds(node *corev1.Node) (*satrapv1.Count, error) {
	url := node.URL("/api/v1/inbounds/count")
	status, resp, err := c.httpClient.Do(http.MethodGet, url, node.Spec.Token, nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errs.New("unexpected", "unexpected response").WithField("nodeID", node.Metadata.ID).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	count := &satrapv1.Count{}
	if err := json.Unmarshal(resp, count); err != nil {
		return nil, errs.New("unmarshal", "unmarshal failed").WithField("nodeID", node.Metadata.ID).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return count, nil
}

func (c *Client) AddInbound(node *corev1.Node, inboundConfig *satrapv1.InboundConfig) error {
	if err := inboundConfig.Validate(); err != nil {
		return fmt.Errorf("validate inbound failed %s/%s: %w", node.Metadata.ID, inboundConfig.Tag, err)
	}
	url := node.URL("/api/v1/inbounds")
	status, resp, err := c.httpClient.Do(http.MethodPost, url, node.Spec.Token, inboundConfig)
	if err != nil {
		return err
	}
	if status == http.StatusConflict {
		return errs.ErrConflict
	}
	if status != http.StatusCreated {
		return errs.New("unexpected", "unexpected response").WithField("nodeID", node.Metadata.ID).WithField("tag", inboundConfig.Tag).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return nil
}

func (c *Client) RemoveInbound(node *corev1.Node, tag string) error {
	path := fmt.Sprintf("/api/v1/inbounds/%s", tag)
	url := node.URL(path)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, node.Spec.Token, nil)
	if err != nil {
		return err
	}
	if status == http.StatusNotFound {
		return errs.ErrNotFound
	}
	if status != http.StatusNoContent {
		return errs.New("unexpected", "unexpected response").WithField("nodeID", node.Metadata.ID).WithField("tag", tag).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return nil
}

func (c *Client) AddUser(node *corev1.Node, tag string, user *satrapv1.InboundUser) error {
	path := fmt.Sprintf("/api/v1/inbounds/%s/users", tag)
	url := node.URL(path)
	status, resp, err := c.httpClient.Do(http.MethodPost, url, node.Spec.Token, user)
	if err != nil {
		return err
	}
	if status == http.StatusConflict {
		return errs.ErrConflict
	}
	if status != http.StatusCreated {
		return errs.New("unexpected", "unexpected response").WithField("nodeID", node.Metadata.ID).WithField("tag", tag).WithField("email", user.Email).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return nil
}

func (c *Client) RemoveUser(node *corev1.Node, tag, email string) error {
	path := fmt.Sprintf("/api/v1/inbounds/%s/users/%s", tag, email)
	url := node.URL(path)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, node.Spec.Token, nil)
	if err != nil {
		return err
	}
	if status == http.StatusNotFound {
		return errs.ErrNotFound
	}
	if status != http.StatusNoContent {
		return errs.New("unexpected", "unexpected response").WithField("nodeID", node.Metadata.ID).WithField("tag", tag).WithField("email", email).WithField("status", strconv.Itoa(status)).WithField("resp", string(resp))
	}
	return nil
}
