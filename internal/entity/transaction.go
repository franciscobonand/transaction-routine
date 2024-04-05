package entity

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

var (
	ErrInvalidAccountID       = errors.New("invalid account id")
	ErrInvalidOperationTypeID = errors.New("invalid operation type id")
	ErrInvalidAmount          = errors.New("invalid amount")
	ErrInvalidEventDate       = errors.New("invalid event date")
	ErrTransactionNotFound    = errors.New("transaction not found")
)

type Transaction struct {
	ID              int             `json:"id"`
	AccountID       int             `json:"account_id"`
	OperationTypeID int             `json:"operation_type_id"`
	Amount          decimal.Decimal `json:"amount"`
	EventDate       time.Time       `json:"event_date"`
}

type TransactionFilter struct {
	ID              *int             `json:"id"`
	AccountID       *int             `json:"account_id"`
	OperationTypeID *int             `json:"operation_type_id"`
	Amount          *decimal.Decimal `json:"amount"`
	EventDate       *time.Time       `json:"event_date"`
}

func (tx *Transaction) Validate(opTypes OperationType) error {
	if tx.AccountID <= 0 {
		return ErrInvalidAccountID
	}
	if tx.Amount.IsZero() {
		return ErrInvalidAmount
	}
	if tx.EventDate.IsZero() {
		return ErrInvalidEventDate
	}
	op, ok := opTypes[tx.OperationTypeID]
	if !ok {
		return ErrInvalidOperationTypeID
	}
	if (op.PositiveAmount && tx.Amount.LessThan(decimal.Zero)) || (!op.PositiveAmount && tx.Amount.GreaterThan(decimal.Zero)) {
		tx.Amount = tx.Amount.Neg()
	}
	return nil
}

func (tx Transaction) ToFilter() TransactionFilter {
	filter := TransactionFilter{}
	if tx.ID != 0 {
		filter.ID = &tx.ID
	}
	if tx.AccountID != 0 {
		filter.AccountID = &tx.AccountID
	}
	if tx.OperationTypeID != 0 {
		filter.OperationTypeID = &tx.OperationTypeID
	}
	if !tx.Amount.IsZero() {
		filter.Amount = &tx.Amount
	}
	if !tx.EventDate.IsZero() {
		filter.EventDate = &tx.EventDate
	}
	return filter
}

func (tx *Transaction) Update(newTx Transaction) {
	if newTx.AccountID != 0 {
		tx.AccountID = newTx.AccountID
	}
	if newTx.OperationTypeID != 0 {
		tx.OperationTypeID = newTx.OperationTypeID
	}
	if !newTx.Amount.IsZero() {
		tx.Amount = newTx.Amount
	}
	if !newTx.EventDate.IsZero() {
		tx.EventDate = newTx.EventDate
	}
}
