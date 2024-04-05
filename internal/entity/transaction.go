package entity

import (
	"errors"
	"time"
)

var (
	ErrInvalidAccountID       = errors.New("invalid account id")
	ErrInvalidOperationTypeID = errors.New("invalid operation type id")
	ErrInvalidAmount          = errors.New("invalid amount")
	ErrInvalidEventDate       = errors.New("invalid event date")
)

type Transaction struct {
	ID              int       `json:"id"`
	AccountID       int       `json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"event_date"`
}

type TransactionFilter struct {
	ID              *int       `json:"id"`
	AccountID       *int       `json:"account_id"`
	OperationTypeID *int       `json:"operation_type_id"`
	Amount          *float64   `json:"amount"`
	EventDate       *time.Time `json:"event_date"`
}

func (tx *Transaction) Validate(opTypes OperationType) error {
	if tx.AccountID <= 0 {
		return ErrInvalidAccountID
	}
	if tx.Amount == 0 {
		return ErrInvalidAmount
	}
	if tx.EventDate.IsZero() {
		return ErrInvalidEventDate
	}
	op, ok := opTypes[tx.OperationTypeID]
	if !ok {
		return ErrInvalidOperationTypeID
	}
	if (op.PositiveAmount && tx.Amount < 0) || (!op.PositiveAmount && tx.Amount > 0) {
		tx.Amount *= -1
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
	if tx.Amount != 0 {
		filter.Amount = &tx.Amount
	}
	if !tx.EventDate.IsZero() {
		filter.EventDate = &tx.EventDate
	}
	return filter
}
