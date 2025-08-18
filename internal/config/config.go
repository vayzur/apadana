package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func Load(configPath string, cfg any) error {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}
