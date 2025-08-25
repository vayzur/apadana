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

func (c *Client) InboundsCount(node *corev1.Node) (*satrapv1.Count, error) {
	url := fmt.Sprintf("%s/api/v1/inbounds/count", node.Address)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, node.Token, nil)
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
	url := fmt.Sprintf("%s/api/v1/inbounds", node.Address)
	status, resp, err := c.httpClient.Do(http.MethodPost, url, node.Token, inboundConfig)
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
	url := fmt.Sprintf("%s/api/v1/inbounds/%s", node.Address, tag)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, node.Token, nil)
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
	url := fmt.Sprintf("%s/api/v1/inbounds/%s/users", node.Address, tag)
	status, resp, err := c.httpClient.Do(http.MethodPost, url, node.Token, user)
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
	url := fmt.Sprintf("%s/api/v1/inbounds/%s/users/%s", node.Address, tag, email)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, node.Token, nil)
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
