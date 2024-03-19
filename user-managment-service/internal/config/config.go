package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env string `envconfig:"ENV"`
	Storage
	Cash
	Broker
	Token
	HTTPServer
}

type HTTPServer struct {
	Address         string        `envconfig:"HTTP_SERVER_ADDRESS"`
	Timeout         time.Duration `envconfig:"HTTP_SERVER_TIMEOUT"`
	IdleTimeout     time.Duration `envconfig:"HTTP_SERVER_IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `envconfig:"HTTP_SERVER_SHUTDOWN_TIMEOUT"`
}

type Storage struct {
	ConnStr string `envconfig:"DB_CONN_STR"`
}

type Cash struct {
	Addr     string `envconfig:"CASH_ADDR"`
	Password string `envconfig:"CASH_PASSWORD"`
	DB       int    `envconfig:"CASH_DB"`
}

type Broker struct {
	ConnStr   string `envconfig:"BROKER_CONN_STR"`
	QueueName string `envconfig:"QUEUE_NAME"`
}

type Token struct {
	JWT struct {
		Secret string        `envconfig:"JWT_TOKEN_SECRET"`
		TTL    time.Duration `envconfig:"JWT_TOKEN_TTL"`
	}
	Refresh struct {
		TTL time.Duration `envconfig:"REFRESH_TOKEN_TTL"`
	}
}

func MustLoad() *Config {
	var cfg Config

	err := godotenv.Load()
	if err != nil {
		log.Panicf("failed to load .env file: %v", err)
	}

	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Panicf("failed to make config: %v", err)
	}

	return &cfg
}
