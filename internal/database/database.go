package database

import (
	"context"
	"errors"
	"fmt"
	"time"
	"transaction-routine/internal/config"
	"transaction-routine/internal/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

const (
	operationTypeTable = "pismo.operation_type"
	accountTable       = "pismo.account"
	transactionTable   = "pismo.transaction"
)

type Service interface {
	Health(ctx context.Context) string
	GetOperationTypes(ctx context.Context) (entity.OperationType, error)
	CreateAccount(ctx context.Context, documentNumber string) error
	GetAccount(ctx context.Context, id int) (*entity.Account, error)
	CreateTransaction(ctx context.Context, t entity.Transaction) error
}

type service struct {
	pool *pgxpool.Pool
	cfg  *config.Config
}

func New(ctx context.Context, cfg *config.Config) (Service, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName,
	)
	pcfg, err := poolConfig(connStr)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, pcfg)
	if err != nil {
		return nil, err
	}

	s := &service{pool: pool, cfg: cfg}
	return s, nil
}

func (s *service) Health(ctx context.Context) string {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return fmt.Sprintf("Error acquiring connection: %v", err)
	}
	defer conn.Release()

	err = conn.Ping(ctx)
	if err != nil {
		return fmt.Sprintf("Error pinging database: %v", err)
	}

	return "It's healthy"
}

func (s *service) GetOperationTypes(ctx context.Context) (entity.OperationType, error) {
	query := fmt.Sprintf("SELECT id, description, positive_amount FROM %s", operationTypeTable)
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ops := make(entity.OperationType)
	for rows.Next() {
		var id int
		var description string
		var positive_amount bool
		err := rows.Scan(&id, &description, &positive_amount)
		if err != nil {
			return nil, err
		}
		ops[id] = entity.Operation{
			Description:    description,
			PositiveAmount: positive_amount,
		}
	}
	return ops, nil
}

func (s *service) CreateAccount(ctx context.Context, documentNumber string) error {
	query := fmt.Sprintf("INSERT INTO %s (document_number) VALUES ($1)", accountTable)
	_, err := s.pool.Exec(
		ctx,
		query,
		documentNumber,
	)
	return err
}

func (s *service) GetAccount(ctx context.Context, id int) (*entity.Account, error) {
	var acc entity.Account
	query := fmt.Sprintf("SELECT id, document_number FROM %s WHERE id = $1", accountTable)
	err := s.pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(&acc.ID, &acc.DocumentNumber)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &acc, err
}

func (s *service) CreateTransaction(ctx context.Context, t entity.Transaction) error {
	query := fmt.Sprintf("INSERT INTO %s (account_id, operation_type_id, amount, event_date) VALUES ($1, $2, $3, $4)", transactionTable)
	_, err := s.pool.Exec(
		ctx,
		query,
		t.AccountID, t.OperationTypeID, t.Amount, t.EventDate,
	)
	return err
}
