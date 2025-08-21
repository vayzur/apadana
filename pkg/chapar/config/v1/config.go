package v1

type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled" yaml:"enabled"`
	CertFile string `mapstructure:"certFile" yaml:"certFile"`
	KeyFile  string `mapstructure:"keyFile" yaml:"keyFile"`
}

type Config struct {
	Address       string    `mapstructure:"address" yaml:"address"`
	Port          uint16    `mapstructure:"port" yaml:"port"`
	Prefork       bool      `mapstructure:"prefork" yaml:"prefork"`
	Token         string    `mapstructure:"token" yaml:"token"`
	TLS           TLSConfig `mapstructure:"tls" yaml:"tls"`
	EtcdEndpoints []string  `mapstructure:"etcd" yaml:"etcd"`
}
