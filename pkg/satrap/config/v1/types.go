package v1

import (
	"time"

	corev1 "github.com/vayzur/apadana/pkg/api/core/v1"
	chaparconfigv1 "github.com/vayzur/apadana/pkg/chapar/config/v1"
	xrayconfigv1 "github.com/vayzur/apadana/pkg/satrap/xray/config/v1"
)

type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled" yaml:"enabled"`
	CertFile string `mapstructure:"certFile" yaml:"certFile"`
	KeyFile  string `mapstructure:"keyFile" yaml:"keyFile"`
}

type SatrapConfig struct {
	NodeID                    string                       `mapstructure:"nodeID" yaml:"nodeID"`
	BindAddress               string                       `mapstructure:"bindAddress" yaml:"bindAddress"`
	Port                      uint16                       `mapstructure:"port" yaml:"port"`
	Prefork                   bool                         `mapstructure:"prefork" yaml:"prefork"`
	Token                     string                       `mapstructure:"token" yaml:"token"`
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
