package config

import (
	"os"
)

type Config struct {
	GRPCPort string `env:"GRPC_PORT"`
	HTTPPort string `env:"HTTP_PORT"`
}

func NewDefaultConfig() *Config {
	return &Config{
		GRPCPort: os.Getenv("GRPC_PORT"),
		HTTPPort: os.Getenv("HTTP_PORT"),
	}
}
