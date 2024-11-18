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
	UpdateTransaction(ctx context.Context, t entity.Transaction) error
}

type transactionService struct {
	cl      clock.Clock
	repo    database.Repository
	opTypes entity.OperationType
	accSvc  AccountService
}

func NewTransactionService(cl clock.Clock, repo database.Repository, opTypes entity.OperationType, accSvc AccountService) TransactionService {
	return &transactionService{cl: cl, repo: repo, opTypes: opTypes, accSvc: accSvc}
}

func (s *transactionService) CreateTransaction(ctx context.Context, t entity.Transaction) error {
	t.EventDate = s.cl.Now()
	if err := t.Validate(s.opTypes); err != nil {
		log.Printf("error validating transaction: %s", err)
		return err
	}

	// lock

	// Buscar a conta referente a transacao
	// validar que a transacao a ser feita "cabe" no limite atual do usuario
	// caso nao caiba, retornar erro
	// 5000 - 70 < 0 => erro

	// open tx

	if err := s.repo.CreateTransaction(ctx, t); err != nil {
		log.Printf("error creating transaction: %s", err)
		// rollback tx
		return err

	}
	// Update account limit

	// acc.AvailCredit += 70
	// s.repo.UpdateAccount(acc)
	// if err, rollback tx
	// commit tx

	// unlock

	return nil
}

func (s *transactionService) UpdateTransaction(ctx context.Context, tx entity.Transaction) error {
	if err := tx.Validate(s.opTypes); err != nil {
		log.Printf("error validating transaction to update: %s", err)
		return err
	}

	currTx, err := s.repo.FindTransactions(ctx, entity.TransactionFilter{ID: &tx.ID})
	if err != nil {
		log.Printf("error getting transaction '%d' to update: %s", tx.ID, err)
		return err
	}
	if len(currTx) == 0 {
		log.Printf("transaction '%d' not found", tx.ID)
		return entity.ErrTransactionNotFound
	}

	newTx := currTx[0]
	newTx.Update(tx)
	if err := s.repo.UpdateTransaction(ctx, newTx); err != nil {
		log.Printf("error updating transaction: %s", err)
		return err
	}
	return nil
}
