package chapar

import (
	"fmt"
	"net/http"

	v1 "github.com/vayzur/apadana/pkg/api/v1"
	"github.com/vayzur/apadana/pkg/errs"
)

func (c *Client) AddInbound(inbound *v1.InboundConfig, node *v1.Node) error {
	if err := inbound.Validate(); err != nil {
		return fmt.Errorf("validate inbound %s/%s: %w", node.Metadata.ID, inbound.Tag, err)
	}

	url := fmt.Sprintf("%s/api/v1/inbounds", node.Address)
	status, resp, err := c.httpClient.Do(http.MethodPost, url, node.Token, inbound)
	if err != nil {
		return fmt.Errorf("add inbound %s/%s: %w", node.Metadata.ID, inbound.Tag, err)
	}
	if status == 409 {
		return errs.ErrConflict
	}
	if status != 201 {
		return fmt.Errorf("add inbound %s/%s: status: %d resp: %s", node.Metadata.ID, inbound.Tag, status, resp)
	}

	return nil
}

func (c *Client) RemoveInbound(node *v1.Node, tag string) error {
	url := fmt.Sprintf("%s/api/v1/inbounds/%s", node.Address, tag)
	status, resp, err := c.httpClient.Do(http.MethodDelete, url, node.Token, nil)
	if err != nil {
		return fmt.Errorf("delete inbound %s/%s: %w", node.Metadata.ID, tag, err)
	}
	if status == 404 {
		return errs.ErrNotFound
	}
	if status != 204 {
		return fmt.Errorf("delete inbound %s/%s: status: %d resp: %s", node.Metadata.ID, tag, status, resp)
	}

	return nil
}
