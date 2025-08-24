package v1

import "time"

type ClusterConfig struct {
	Enabled bool   `mapstructure:"enabled" yaml:"enabled"`
	Server  string `mapstructure:"server" yaml:"server"`
	Token   string `mapstructure:"token" yaml:"token"`
}

type SpasakaConfig struct {
	Cluster                ClusterConfig `mapstructure:"cluster" yaml:"cluster"`
	EtcdEndpoints          []string      `mapstructure:"etcd" yaml:"etcd"`
	ConcurrentNodeSyncs    int           `mapstructure:"concurrentNodeSyncs" yaml:"concurrentNodeSyncs"`
	NodeMonitorPeriod      time.Duration `mapstructure:"nodeMonitorPeriod" yaml:"nodeMonitorPeriod"`
	NodeMonitorGracePeriod time.Duration `mapstructure:"nodeMonitorGracePeriod" yaml:"nodeMonitorGracePeriod"`
}
