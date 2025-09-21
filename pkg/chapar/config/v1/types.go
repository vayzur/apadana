package v1

import etcdconfigv1 "github.com/vayzur/apadana/pkg/chapar/storage/etcd/config/v1"

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

type ChaparConfig struct {
	Address string                  `mapstructure:"address" yaml:"address"`
	Port    uint16                  `mapstructure:"port" yaml:"port"`
	Prefork bool                    `mapstructure:"prefork" yaml:"prefork"`
	Token   string                  `mapstructure:"token" yaml:"token"`
	TLS     TLSConfig               `mapstructure:"tls" yaml:"tls"`
	Etcd    etcdconfigv1.EtcdConfig `mapstructure:"etcd" yaml:"etcd"`
}
