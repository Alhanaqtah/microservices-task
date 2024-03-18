package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"user-managment-service/internal/models"
	"user-managment-service/internal/storage/storage"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type Storage interface {
	User(ctx context.Context, username string) (*models.User, error)
	CreateNewUser(ctx context.Context, username string, passHash []byte) (string, error)
}

type Cash interface {
}

type Broker interface {
}

type Service struct {
	log     *slog.Logger
	storage Storage
	cash    Cash
	broker  Broker
}

func New(log *slog.Logger, storage Storage, cash Cash, broker Broker) *Service {
	return &Service{
		log:     log,
		storage: storage,
		cash:    cash,
		broker:  broker,
	}
}

func (s *Service) SignUp(username string, password string) (string, error) {
	const op = "service.auth.SignUp"

	_ = s.log.With(slog.String("op", op))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	u, err := s.storage.User(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrUserNotFound) {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if u != nil {
		return "", fmt.Errorf("%s: %w", op, storage.ErrUserExists)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	uuid, err := s.storage.CreateNewUser(ctx, username, passHash)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return uuid, nil
}
