package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	nullBytes  = []byte("null")
	emptyObj   = []byte("{}")
	emptyArray = []byte("[]")
)

type InboundConfig struct {
	Tag            string          `json:"tag"`
	Protocol       string          `json:"protocol"`
	Port           uint16          `json:"port"`
	Listen         json.RawMessage `json:"listen"`
	Settings       json.RawMessage `json:"settings"`
	Allocate       json.RawMessage `json:"allocate"`
	StreamSettings json.RawMessage `json:"streamSettings"`
	Sniffing       json.RawMessage `json:"sniffing"`
}

type Inbound struct {
	Metadata Metadata      `json:"metadata"`
	Config   InboundConfig `json:"config"`
}

func (c *InboundConfig) Validate() error {
	if c.Tag == "" {
		return errors.New("tag cannot be empty")
	}
	if c.Protocol == "" {
		return errors.New("protocol cannot be empty")
	}
	if c.Port == 0 {
		return errors.New("port cannot be zero")
	}

	fields := map[string]json.RawMessage{
		"listen":         c.Listen,
		"settings":       c.Settings,
		"allocate":       c.Allocate,
		"streamSettings": c.StreamSettings,
		"sniffing":       c.Sniffing,
	}

	for name, value := range fields {
		if isRawMessageEmpty(value) {
			return fmt.Errorf("%s cannot be empty", name)
		}
	}

	return nil
}

func isRawMessageEmpty(data json.RawMessage) bool {
	trimmed := bytes.TrimSpace(data)
	return len(trimmed) == 0 ||
		bytes.Equal(trimmed, nullBytes) ||
		bytes.Equal(trimmed, emptyObj) ||
		bytes.Equal(trimmed, emptyArray)
}
