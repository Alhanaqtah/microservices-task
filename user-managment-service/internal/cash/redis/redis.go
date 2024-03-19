package redis

import (
	"context"
	"fmt"
	"user-managment-service/internal/config"

	"github.com/redis/go-redis/v9"
)

type Cash struct {
	client *redis.Client
}

func New(cfg config.Cash) (*Cash, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	err = client.SAdd(ctx, "blacklist", "").Err()
	if err != nil {
		return nil, err
	}

	return &Cash{client: client}, nil
}

func (c *Cash) AddToBlaclist(ctx context.Context, token string) error {
	const op = "SearchInBlacklist"

	err := c.client.SAdd(ctx, "blacklist", token).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Cash) SearchInBlacklist(ctx context.Context, token string) (bool, error) {
	const op = "SearchInBlacklist"

	found, err := c.client.SIsMember(ctx, "blacklist", token).Result()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return found, nil
}
