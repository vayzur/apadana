package v1

type XrayConfig struct {
	Address string `mapstructure:"address" yaml:"address"`
	Port    uint16 `mapstructure:"port" yaml:"port"`
}
