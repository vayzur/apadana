package v1

import (
	"time"

	"github.com/google/uuid"
	corev1 "github.com/vayzur/apadana/pkg/apis/core/v1"
	chaparconfigv1 "github.com/vayzur/apadana/pkg/chapar/config/v1"
	xrayconfigv1 "github.com/vayzur/apadana/pkg/satrap/xray/config/v1"
)

type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled" yaml:"enabled"`
	CertFile string `mapstructure:"certFile" yaml:"certFile"`
	KeyFile  string `mapstructure:"keyFile" yaml:"keyFile"`
}

type SatrapConfig struct {
	Name                      string                       `mapstructure:"name" yaml:"name"`
	Address                   string                       `mapstructure:"address" yaml:"address"`
	Port                      uint16                       `mapstructure:"port" yaml:"port"`
	Prefork                   bool                         `mapstructure:"prefork" yaml:"prefork"`
	Token                     string                       `mapstructure:"token" yaml:"token"`
	RegisterNode              bool                         `mapstructure:"registerNode" yaml:"registerNode"`
	Addresses                 []corev1.NodeAddress         `mapstructure:"addresses" yaml:"addresses"`
	TLS                       TLSConfig                    `mapstructure:"tls" yaml:"tls"`
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

func (c *SatrapConfig) GetToken() string {
	// If user explicitly set a token, always use it (both standalone and cluster mode)
	if c.Token != "" {
		return c.Token
	}

	// No token provided by user
	if c.Cluster.Enabled && c.RegisterNode {
		// Cluster mode: generate random token for auto-registration security
		token := uuid.NewString()
		return token
	}

	return ""
}
