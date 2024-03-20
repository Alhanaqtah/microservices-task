package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"user-managment-service/internal/models"
	"user-managment-service/internal/storage/storage"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Storage interface {
	UserByUUID(ctx context.Context, uuid string) (*models.User, error)
}

type Service struct {
	log     *slog.Logger
	storage Storage
}

func New(log *slog.Logger, storage Storage) *Service {
	return &Service{
		log:     log,
		storage: storage,
	}
}

func (s *Service) UserByUUID(uuid string) (*models.User, error) {
	const op = "service.user.UserByUUID"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user, err := s.storage.UserByUUID(ctx, uuid)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
