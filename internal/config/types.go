package config

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
