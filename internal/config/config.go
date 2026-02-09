package config

import (
	"os"
	"time"
)

const (
	defaultPort         = "8080"
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 10 * time.Second
)

type HttpServerCfg struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func LoadHttpServerConfig() *HttpServerCfg {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = defaultPort
	}

	return &HttpServerCfg{
		Port:         port,
		ReadTimeout:  getEnvAsDuration("HTTP_READ_TIMEOUT", defaultReadTimeout),
		WriteTimeout: getEnvAsDuration("HTTP_WRITE_TIMEOUT", defaultWriteTimeout),
	}
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}

	duration, err := time.ParseDuration(valStr)
	if err != nil {
		return fallback
	}

	return duration
}
