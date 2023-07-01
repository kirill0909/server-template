package redis

import (
	"auth-svc/config"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	opts := &redis.Options{}
	if cfg.Redis.UseCertificates {
		certs := make([]tls.Certificate, 0, 0)
		if cfg.Redis.CertificatesPaths.Cert != "" && cfg.Redis.CertificatesPaths.Key != "" {
			cert, err := tls.LoadX509KeyPair(cfg.Redis.CertificatesPaths.Cert, cfg.Redis.CertificatesPaths.Key)
			if err != nil {
				return nil, errors.Wrapf(
					err,
					"certPath: %v, keyPath: %v",
					cfg.Redis.CertificatesPaths.Cert,
					cfg.Redis.CertificatesPaths.Key,
				)
			}
			certs = append(certs, cert)
		}
		caCert, err := os.ReadFile(cfg.Redis.CertificatesPaths.Ca)
		if err != nil {
			return nil, errors.Wrapf(err, "ca load path: %v", cfg.Redis.CertificatesPaths.Ca)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		opts = &redis.Options{
			Addr:         fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
			MinIdleConns: cfg.Redis.MinIdleConns,
			PoolSize:     cfg.Redis.PoolSize,
			PoolTimeout:  time.Duration(cfg.Redis.PoolTimeout) * time.Second,
			Password:     cfg.Redis.Password,
			DB:           cfg.Redis.DB,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: cfg.Redis.InsecureSkipVerify,
				Certificates:       certs,
				RootCAs:            caCertPool,
			},
		}
	} else {
		opts = &redis.Options{
			Addr:         fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
			MinIdleConns: cfg.Redis.MinIdleConns,
			PoolSize:     cfg.Redis.PoolSize,
			PoolTimeout:  time.Duration(cfg.Redis.PoolTimeout) * time.Second,
			Password:     cfg.Redis.Password,
			DB:           cfg.Redis.DB,
		}
	}

	client := redis.NewClient(opts)
	result := client.Ping(context.Background())
	if result.Err() != nil {
		return nil, result.Err()
	}

	return client, nil
}
