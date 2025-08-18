package config

type ChaparConfig struct {
	Address       string    `mapstructure:"address" yaml:"address"`
	Port          uint16    `mapstructure:"port" yaml:"port"`
	Prefork       bool      `mapstructure:"prefork" yaml:"prefork"`
	Token         string    `mapstructure:"token" yaml:"token"`
	TLS           TLSConfig `mapstructure:"tls" yaml:"tls"`
	EtcdEndpoints []string  `mapstructure:"etcd" yaml:"etcd"`
}
