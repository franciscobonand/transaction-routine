package service

import (
	"context"
	"log"
	"transaction-routine/internal/database"
)

type HealthService interface {
	HealthCheck(ctx context.Context) bool
}

type healthService struct {
	repo database.Repository
}

func NewHealthService(repo database.Repository) HealthService {
	return &healthService{repo: repo}
}

func (s *healthService) HealthCheck(ctx context.Context) bool {
	if err := s.repo.Health(ctx); err != nil {
		log.Panicf("error checking database health: %s", err)
		return false
	}
	return true
}
