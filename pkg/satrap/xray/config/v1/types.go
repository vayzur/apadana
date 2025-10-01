package v1

import "time"

type XrayConfig struct {
	Address               string        `mapstructure:"address" yaml:"address"`
	Port                  uint16        `mapstructure:"port" yaml:"port"`
	RuntimeRequestTimeout time.Duration `mapstructure:"runtimeRequestTimeout" yaml:"runtimeRequestTimeout"`
}
