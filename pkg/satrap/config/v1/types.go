package v1

import (
	"os"
	"time"

	corev1 "github.com/vayzur/apadana/pkg/apis/core/v1"
	chaparconfigv1 "github.com/vayzur/apadana/pkg/chapar/config/v1"
	xrayconfigv1 "github.com/vayzur/apadana/pkg/satrap/xray/config/v1"
)

type SatrapConfig struct {
	Name                      string                       `mapstructure:"name" yaml:"name"`
	RegisterNode              bool                         `mapstructure:"registerNode" yaml:"registerNode"`
	Addresses                 []corev1.NodeAddress         `mapstructure:"addresses" yaml:"addresses"`
	Xray                      xrayconfigv1.XrayConfig      `mapstructure:"xray" yaml:"xray"`
	Cluster                   chaparconfigv1.ClusterConfig `mapstructure:"cluster" yaml:"cluster"`
	NodeStatusUpdateFrequency time.Duration                `mapstructure:"nodeStatusUpdateFrequency" yaml:"nodeStatusUpdateFrequency"`
	SyncFrequency             time.Duration                `mapstructure:"syncFrequency" yaml:"syncFrequency"`
	ConcurrentInboundSyncs    uint32                       `mapstructure:"concurrentInboundSyncs" yaml:"concurrentInboundSyncs"`
	ConcurrentInboundGCSyncs  uint32                       `mapstructure:"concurrentInboundGCSyncs" yaml:"concurrentInboundGCSyncs"`
	ConcurrentUserSyncs       uint32                       `mapstructure:"concurrentUserSyncs" yaml:"concurrentUserSyncs"`
	ConcurrentUserGCSyncs     uint32                       `mapstructure:"concurrentUserGCSyncs" yaml:"concurrentUserGCSyncs"`
	MaxInbounds               uint32                       `mapstructure:"maxInbounds" yaml:"maxInbounds"`
}

func (c *SatrapConfig) GetName() string {
	if c.Name != "" {
		return c.Name
	}
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}
