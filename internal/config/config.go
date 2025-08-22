package config

import (
	"log"
	"my_blog_backend/pkg/e"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return e.Wrap("error loading .env file", err)
	}

	return nil
}

type HttpServer struct {
	Port         string        `mapstructure:"HTTP_PORT"`
	ReadTimeout  time.Duration `mapstructure:"HTTP_READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"HTTP_WRITE_TIMEOUT"`
}

func LoadHttpServerConfig() HttpServer {
	v := viper.New()

	// Берём переменные из окружения (godotenv уже их загрузил)
	v.AutomaticEnv()

	var cfg HttpServer
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("failed to unmarshal HttpServer config: %v", err)
	}

	return cfg
}
