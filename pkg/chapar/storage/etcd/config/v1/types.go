package v1

type EtcdTLSConfig struct {
	Enabled  bool   `mapstructure:"enabled" yaml:"enabled"`
	CAFile   string `mapstructure:"caFile" yaml:"caFile"`
	CertFile string `mapstructure:"certFile" yaml:"certFile"`
	KeyFile  string `mapstructure:"keyFile" yaml:"keyFile"`
}

type EtcdConfig struct {
	Servers []string      `mapstructure:"servers" yaml:"servers"`
	TLS     EtcdTLSConfig `mapstructure:"tls" yaml:"tls"`
}
