package database

import (
	"context"
	"fmt"
	"os"
	"time"
	"transaction-routine/internal/entity"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Health(ctx context.Context) map[string]string
	GetOperationTypes(ctx context.Context) (entity.OperationType, error)
	CreateAccount(ctx context.Context, documentNumber string) error
	GetAccount(ctx context.Context, id int) (entity.Account, error)
	CreateTransaction(ctx context.Context, t entity.Transaction) error
}

type service struct {
	pool *pgxpool.Pool
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New(ctx context.Context) (Service, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	cfg, err := config(connStr)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	s := &service{pool: pool}
	return s, nil
}

func (s *service) Health(ctx context.Context) map[string]string {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return map[string]string{
			"message": fmt.Sprintf("Error acquiring connection: %v", err),
		}
	}
	defer conn.Release()

	err = conn.Ping(ctx)
	if err != nil {
		return map[string]string{
			"message": fmt.Sprintf("Error pinging database: %v", err),
		}
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) GetOperationTypes(ctx context.Context) (entity.OperationType, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, description FROM operation_types")
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
	_, err := s.pool.Exec(
		ctx,
		"INSERT INTO accounts (document_number) VALUES ($1)",
		documentNumber,
	)
	return err
}

func (s *service) GetAccount(ctx context.Context, id int) (entity.Account, error) {
	var acc entity.Account
	err := s.pool.QueryRow(
		ctx,
		"SELECT id, document_number FROM accounts WHERE id = $1",
		id,
	).Scan(&acc)
	return acc, err
}

func (s *service) CreateTransaction(ctx context.Context, t entity.Transaction) error {
	_, err := s.pool.Exec(
		ctx,
		"INSERT INTO transactions (account_id, operation_type_id, amount, event_date) VALUES ($1, $2, $3, $4)",
		t.AccountID, t.OperationTypeID, t.Amount, t.EventDate,
	)
	return err
}
