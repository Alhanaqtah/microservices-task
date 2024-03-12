package redis

import (
	"user-managment-service/internal/config"

	"github.com/redis/go-redis/v9"
)

type Cash struct {
	client *redis.Client
}

func New(cfg config.Cash) *Cash {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &Cash{client: client}
}
