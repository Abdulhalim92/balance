package service

import (
	"balance/internal/repository"
	"balance/logging"
	"github.com/go-redis/redis"
)

type Service struct {
	Repository *repository.Repository
	Logger     *logging.Logger
	Redis      *redis.Client
}

func NewService(repository *repository.Repository, logger *logging.Logger, redis *redis.Client) *Service {
	return &Service{
		Logger:     logger,
		Repository: repository,
		Redis:      redis,
	}
}
