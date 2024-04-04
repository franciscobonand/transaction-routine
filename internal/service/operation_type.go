//go:generate mockgen -destination=./../../tests/mocks/mock_operation_type.go -package=mocks -source=operation_type.go
package service

import (
	"context"
	"log"
	"transaction-routine/internal/database"
	"transaction-routine/internal/entity"
)

type OpTypeService interface {
	GetAllOperationTypes(ctx context.Context) (entity.OperationType, error)
	CreateOperationType(ctx context.Context, op entity.Operation) error
}

type opTypeService struct {
	repo    database.Repository
	opTypes entity.OperationType
}

func NewOpTypeService(repo database.Repository, opTypes entity.OperationType) OpTypeService {
	return &opTypeService{repo: repo, opTypes: opTypes}
}

func (s *opTypeService) GetAllOperationTypes(ctx context.Context) (entity.OperationType, error) {
	return s.repo.FindOperationType(ctx)
}

func (s *opTypeService) CreateOperationType(ctx context.Context, op entity.Operation) error {
	if err := s.repo.CreateOperationType(ctx, op); err != nil {
		log.Printf("error creating operation type: %s", err)
		return err
	}
	return nil
}

func (s *opTypeService) RefreshOperationTypes(ctx context.Context) error {
	opTypes, err := s.repo.FindOperationType(ctx)
	if err != nil {
		log.Printf("error getting operation types: %s", err)
		return err
	}
	s.opTypes = opTypes
	return nil
}
