//go:generate mockgen -destination=./../../tests/mocks/mock_transaction.go -package=mocks -source=transaction.go
package service

import (
	"context"
	"log"
	"transaction-routine/internal/clock"
	"transaction-routine/internal/database"
	"transaction-routine/internal/entity"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, t entity.Transaction) error
}

type transactionService struct {
	cl      clock.Clock
	repo    database.Repository
	opTypes entity.OperationType
}

func NewTransactionService(cl clock.Clock, repo database.Repository, opTypes entity.OperationType) TransactionService {
	return &transactionService{cl: cl, repo: repo, opTypes: opTypes}
}

func (s *transactionService) CreateTransaction(ctx context.Context, t entity.Transaction) error {
	t.EventDate = s.cl.Now()
	if err := t.Validate(s.opTypes); err != nil {
		log.Printf("error validating transaction: %s", err)
		return err
	}
	if err := s.repo.CreateTransaction(ctx, t); err != nil {
		log.Printf("error creating transaction: %s", err)
		return err
	}
	return nil
}
