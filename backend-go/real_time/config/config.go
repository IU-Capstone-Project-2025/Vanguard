package config

import (
	"os"
)

// ServiceConfig is a structure containing all loaded variables from environment
type ServiceConfig struct {
	Host string       // server host
	Port string       // server port
	MQ   RabbitConfig // Message broker configs
	JWT  JWTConfig    // Jwt configs
}

// RabbitConfig is a structure containing environment variables for RabbitMQ setup
type RabbitConfig struct {
	User     string
	Password string
	Host     string
	Port     string
}

// JWTConfig is a structure containing environment variables related to JWT
type JWTConfig struct {
	SecretKey string // jwt secret key
}

// config stores once parsed env variables
var config *ServiceConfig

// LoadConfig is a singleton functions, that returns parsed config.
// If the function have not been called, then it parses data from environment and stores in `config` variable.
// Otherwise, just returns already parsed config
func LoadConfig() *ServiceConfig {
	if config != nil {
		return config
	}

	cfg := &ServiceConfig{
		Host: os.Getenv("REALTIME_SERVICE_HOST"),
		Port: os.Getenv("REALTIME_SERVICE_PORT"),
		MQ: RabbitConfig{
			User:     os.Getenv("RABBITMQ_USER"),
			Password: os.Getenv("RABBITMQ_PASSWORD"),
			Host:     os.Getenv("RABBITMQ_HOST"),
			Port:     os.Getenv("RABBITMQ_PORT"),
		},
		JWT: JWTConfig{SecretKey: os.Getenv("JWT_SECRET_KEY")},
	}

	config = cfg

	return cfg
}
