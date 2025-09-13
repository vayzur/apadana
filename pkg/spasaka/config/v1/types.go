package v1

import (
	"time"

	chaparconfigv1 "github.com/vayzur/apadana/pkg/chapar/config/v1"
	etcdconfigv1 "github.com/vayzur/apadana/pkg/chapar/storage/etcd/config/v1"
)

type SpasakaConfig struct {
	Cluster                chaparconfigv1.ClusterConfig `mapstructure:"cluster" yaml:"cluster"`
	Etcd                   etcdconfigv1.EtcdConfig      `mapstructure:"etcd" yaml:"etcd"`
	ConcurrentNodeSyncs    int                          `mapstructure:"concurrentNodeSyncs" yaml:"concurrentNodeSyncs"`
	NodeMonitorPeriod      time.Duration                `mapstructure:"nodeMonitorPeriod" yaml:"nodeMonitorPeriod"`
	NodeMonitorGracePeriod time.Duration                `mapstructure:"nodeMonitorGracePeriod" yaml:"nodeMonitorGracePeriod"`
}
