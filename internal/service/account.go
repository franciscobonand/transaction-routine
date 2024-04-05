//go:generate mockgen -destination=./../../tests/mocks/mock_account.go -package=mocks -source=account.go
package service

import (
	"context"
	"log"
	"transaction-routine/internal/database"
	"transaction-routine/internal/entity"

	"github.com/shopspring/decimal"
)

type AccountService interface {
	GetAccountByID(ctx context.Context, id int) (*entity.Account, error)
	CreateAccount(ctx context.Context, acc entity.Account) error
	GetAccountBalance(ctx context.Context, id int) (decimal.Decimal, error)
}

type accountService struct {
	repo database.Repository
}

func NewAccountService(repo database.Repository) AccountService {
	return &accountService{repo: repo}
}

func (s *accountService) GetAccountByID(ctx context.Context, id int) (*entity.Account, error) {
	accs, err := s.repo.FindAccounts(ctx, entity.AccountFilter{ID: &id})
	if err != nil {
		log.Printf("error getting account %d: %s", err, err)
		return nil, err
	}
	if len(accs) == 0 {
		return nil, nil
	}
	return &accs[0], nil
}

func (s *accountService) CreateAccount(ctx context.Context, acc entity.Account) error {
	if acc.DocumentNumber == "" {
		return entity.ErrMissingDocumentNumber
	}
	if err := s.repo.CreateAccount(ctx, acc); err != nil {
		log.Printf("error creating account: %s", err)
		return err
	}
	return nil
}

func (s *accountService) GetAccountBalance(ctx context.Context, id int) (decimal.Decimal, error) {
	txs, err := s.repo.FindTransactions(ctx, entity.TransactionFilter{AccountID: &id})
	if err != nil {
		log.Printf("error getting transactions to calculate balance: %s", err)
		return decimal.Zero, err
	}
	balance := decimal.Zero
	for _, tx := range txs {
		balance = balance.Add(tx.Amount)
	}
	return balance, nil
}
