package auth

import (
	"log/slog"
)

type Storage interface {
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
