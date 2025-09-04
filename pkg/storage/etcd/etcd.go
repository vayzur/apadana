package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"time"

	etcdconfigv1 "github.com/vayzur/apadana/pkg/storage/etcd/config/v1"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewClient(cfg *etcdconfigv1.EtcdConfig, ctx context.Context) (*clientv3.Client, error) {
	var tlsConfig *tls.Config

	if cfg.TLS.Enabled {
		caCert, err := os.ReadFile(cfg.TLS.CAFile)
		if err != nil {
			return nil, err
		}
		caPool := x509.NewCertPool()
		if !caPool.AppendCertsFromPEM(caCert) {
			return nil, err
		}

		cert, err := tls.LoadX509KeyPair(cfg.TLS.CertFile, cfg.TLS.KeyFile)
		if err != nil {
			return nil, err
		}

		tlsConfig = &tls.Config{
			RootCAs:      caPool,
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS13,
		}
	}

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Servers,
		DialTimeout: 5 * time.Second,
		Context:     ctx,
		TLS:         tlsConfig,
	})
	if err != nil {
		return nil, err
	}

	return etcdClient, nil
}
