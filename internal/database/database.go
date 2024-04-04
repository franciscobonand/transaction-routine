//go:generate mockgen -destination=./../../tests/mocks/mock_repository.go -package=mocks -source=database.go
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

type Repository interface {
	Health(ctx context.Context) error
	CreateOperationType(ctx context.Context, op entity.Operation) error
	FindOperationType(ctx context.Context) (entity.OperationType, error)
	CreateAccount(ctx context.Context, acc entity.Account) error
	FindAccounts(ctx context.Context, filter entity.AccountFilter) ([]entity.Account, error)
	CreateTransaction(ctx context.Context, tx entity.Transaction) error
	FindTransactions(ctx context.Context, filter entity.TransactionFilter) ([]entity.Transaction, error)
}

type repo struct {
	pool *pgxpool.Pool
	cfg  *config.Config
}

func New(ctx context.Context, cfg *config.Config) (Repository, error) {
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

	s := &repo{pool: pool, cfg: cfg}
	return s, nil
}

func (r *repo) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("error acquiring connection: %v", err)
	}
	defer conn.Release()

	err = conn.Ping(ctx)
	if err != nil {
		return fmt.Errorf("error pinging database: %v", err)
	}

	return nil
}

func (r *repo) CreateOperationType(ctx context.Context, op entity.Operation) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			description,
			positive_amount
		) VALUES ($1, $2)`,
		operationTypeTable,
	)
	_, err := r.pool.Exec(
		ctx,
		query,
		op.Description,
		op.PositiveAmount,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) FindOperationType(ctx context.Context) (entity.OperationType, error) {
	query := fmt.Sprintf("SELECT id, description, positive_amount FROM %s", operationTypeTable)
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ops := make(entity.OperationType)
	for rows.Next() {
		var id int
		var op entity.Operation
		err := rows.Scan(&id, &op.Description, &op.PositiveAmount)
		if err != nil {
			return nil, err
		}
		ops[id] = &op
	}
	return ops, nil
}

func (r *repo) CreateAccount(ctx context.Context, acc entity.Account) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			document_number
		) VALUES ($1)`,
		accountTable,
	)
	_, err := r.pool.Exec(
		ctx,
		query,
		acc.DocumentNumber,
	)
	return err
}

func (r *repo) FindAccounts(ctx context.Context, filter entity.AccountFilter) ([]entity.Account, error) {
	query := fmt.Sprintf(`
		SELECT
			id,
			document_number
		FROM %s
		WHERE
			(id = COALESCE($1, id))
			AND (document_number = COALESCE($2, document_number))
		`,
		accountTable,
	)
	rows, err := r.pool.Query(
		ctx,
		query,
		filter.ID,
		filter.DocumentNumber,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	accs := make([]entity.Account, 0)
	for rows.Next() {
		var acc entity.Account
		err := rows.Scan(&acc.ID, &acc.DocumentNumber)
		if err != nil {
			return nil, err
		}
		accs = append(accs, acc)
	}
	return accs, err
}

func (r *repo) CreateTransaction(ctx context.Context, tx entity.Transaction) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			account_id,
			operation_type_id,
			amount,
			event_date
		) VALUES ($1, $2, $3, $4)`,
		transactionTable,
	)
	_, err := r.pool.Exec(
		ctx,
		query,
		tx.AccountID, tx.OperationTypeID, tx.Amount, tx.EventDate,
	)
	return err
}

func (r *repo) FindTransactions(ctx context.Context, filter entity.TransactionFilter) ([]entity.Transaction, error) {
	query := fmt.Sprintf(`
		SELECT
			id,
			account_id,
			operation_type_id,
			amount,
			event_date
		FROM %s
		WHERE
			(id = COALESCE($1, id))
			AND (account_id = COALESCE($2, account_id))
			AND (operation_type_id = COALESCE($3, operation_type_id))
			AND (amount = COALESCE($4, amount))
			AND (event_date = COALESCE($5, event_date))
		`,
		transactionTable,
	)
	rows, err := r.pool.Query(
		ctx,
		query,
		filter.ID,
		filter.AccountID,
		filter.OperationTypeID,
		filter.Amount,
		filter.EventDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	txs := make([]entity.Transaction, 0)
	for rows.Next() {
		var tx entity.Transaction
		err := rows.Scan(&tx.ID, &tx.AccountID, &tx.OperationTypeID, &tx.Amount, &tx.EventDate)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, err
}
