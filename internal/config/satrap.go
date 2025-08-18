package config

import "time"

type XrayConfig struct {
	Address string `mapstructure:"address" yaml:"address"`
	Port    uint16 `mapstructure:"port" yaml:"port"`
}

type SatrapConfig struct {
	NodeID                    string        `mapstructure:"nodeID" yaml:"nodeID"`
	Address                   string        `mapstructure:"address" yaml:"address"`
	Port                      uint16        `mapstructure:"port" yaml:"port"`
	Prefork                   bool          `mapstructure:"prefork" yaml:"prefork"`
	Token                     string        `mapstructure:"token" yaml:"token"`
	TLS                       TLSConfig     `mapstructure:"tls" yaml:"tls"`
	Xray                      XrayConfig    `mapstructure:"xray" yaml:"xray"`
	Cluster                   ClusterConfig `mapstructure:"cluster" yaml:"cluster"`
	NodeStatusUpdateFrequency time.Duration `mapstructure:"nodeStatusUpdateFrequency" yaml:"nodeStatusUpdateFrequency"`
}
