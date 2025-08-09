package config

import (
	"github.com/joho/godotenv"
	"my_blog_backend/pkg/e"
)

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return e.Wrap("error loading .env file", err)
	}

	return nil
}
