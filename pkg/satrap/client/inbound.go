package satrap

import (
	"encoding/json"
	"fmt"
	"net/http"

	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
	satrapv1 "github.com/vayzur/apadana/pkg/api/satrap/v1"
	"github.com/vayzur/apadana/pkg/errs"
)

func (c *Client) InboundsCount(node *corev1.Node) (*satrapv1.Count, error) {
	url := fmt.Sprintf("%s/api/v1/inbounds/count", node.Address)
	status, resp, err := c.httpClient.Do(http.MethodGet, url, node.Token, nil)
	if err != nil {
		return nil, fmt.Errorf("get inbounds count %s: %w", node.Metadata.ID, err)
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("get inbounds count %s: status: %d resp: %s", node.Metadata.ID, status, resp)
	}
	count := &satrapv1.Count{}
	if err := json.Unmarshal(resp, count); err != nil {
		return nil, fmt.Errorf("unmarshal inbounds count %s: status: %d resp: %s", node.Metadata.ID, status, resp)
	}
	return count, nil
}

func (c *Client) AddInbound(inbound *satrapv1.InboundConfig, node *corev1.Node) error {
	if err := inbound.Validate(); err != nil {
		return fmt.Errorf("validate inbound %s/%s: %w", node.Metadata.ID, inbound.Tag, err)
	}
	url := fmt.Sprintf("%s/api/v1/inbounds", node.Address)
	status, resp, err := c.httpClient.Do(http.MethodPost, url, node.Token, inbound)
	if err != nil {
		return fmt.Errorf("add inbound %s/%s: %w", node.Metadata.ID, inbound.Tag, err)
	}
	if status == http.StatusConflict {
		return errs.ErrConflict
	}
	if status != http.StatusCreated {
		return fmt.Errorf("add inbound %s/%s: status: %d resp: %s", node.Metadata.ID, inbound.Tag, status, resp)
	}
	return nil
}

func (c *Client) RemoveInbound(node *corev1.Node, tag string) error {
	url := fmt.Sprintf("%s/api/v1/inbounds/%s", node.Address, tag)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, node.Token, nil)
	if err != nil {
		return fmt.Errorf("delete inbound %s/%s: %w", node.Metadata.ID, tag, err)
	}
	if status == http.StatusNotFound {
		return errs.ErrNotFound
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("delete inbound %s/%s: status: %d resp: %s", node.Metadata.ID, tag, status, resp)
	}
	return nil
}

func (c *Client) AddUser(node *corev1.Node, tag string, user *satrapv1.InboundUser) error {
	url := fmt.Sprintf("%s/api/v1/inbounds/%s/users", node.Address, tag)
	status, resp, err := c.httpClient.Do(http.MethodPost, url, node.Token, user)
	if err != nil {
		return fmt.Errorf("add user %s/%s: %w", node.Metadata.ID, tag, err)
	}
	if status == http.StatusConflict {
		return errs.ErrConflict
	}
	if status != http.StatusCreated {
		return fmt.Errorf("add user %s/%s: status: %d resp: %s", node.Metadata.ID, tag, status, resp)
	}
	return nil
}

func (c *Client) RemoveUser(node *corev1.Node, tag, email string) error {
	url := fmt.Sprintf("%s/api/v1/inbounds/%s/users/%s", node.Address, tag, email)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, node.Token, nil)
	if err != nil {
		return fmt.Errorf("delete user %s/%s/%s: %w", node.Metadata.ID, tag, email, err)
	}
	if status == http.StatusNotFound {
		return errs.ErrNotFound
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("delete user %s/%s/%s: status: %d resp: %s", node.Metadata.ID, tag, email, status, resp)
	}
	return nil
}
