package v1

import "time"

type SpasakaConfig struct {
	EtcdEndpoints          []string      `mapstructure:"etcd" yaml:"etcd"`
	NodeMonitorPeriod      time.Duration `mapstructure:"nodeMonitorPeriod" yaml:"nodeMonitorPeriod"`
	NodeMonitorGracePeriod time.Duration `mapstructure:"nodeMonitorGracePeriod" yaml:"nodeMonitorGracePeriod"`
}
