package v1

import "time"

type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled" yaml:"enabled"`
	CertFile string `mapstructure:"certFile" yaml:"certFile"`
	KeyFile  string `mapstructure:"keyFile" yaml:"keyFile"`
}
type ClusterConfig struct {
	Enabled bool   `mapstructure:"enabled" yaml:"enabled"`
	Server  string `mapstructure:"server" yaml:"server"`
	Token   string `mapstructure:"token" yaml:"token"`
}

type XrayConfig struct {
	Address string `mapstructure:"address" yaml:"address"`
	Port    uint16 `mapstructure:"port" yaml:"port"`
}

type Config struct {
	NodeID                    string        `mapstructure:"nodeID" yaml:"nodeID"`
	Address                   string        `mapstructure:"address" yaml:"address"`
	Port                      uint16        `mapstructure:"port" yaml:"port"`
	Prefork                   bool          `mapstructure:"prefork" yaml:"prefork"`
	Token                     string        `mapstructure:"token" yaml:"token"`
	TLS                       TLSConfig     `mapstructure:"tls" yaml:"tls"`
	Xray                      XrayConfig    `mapstructure:"xray" yaml:"xray"`
	Cluster                   ClusterConfig `mapstructure:"cluster" yaml:"cluster"`
	NodeStatusUpdateFrequency time.Duration `mapstructure:"nodeStatusUpdateFrequency" yaml:"nodeStatusUpdateFrequency"`
	InboundTTLCheckPeriod     time.Duration `mapstructure:"inboundTTLCheckPeriod" yaml:"inboundTTLCheckPeriod"`
}
