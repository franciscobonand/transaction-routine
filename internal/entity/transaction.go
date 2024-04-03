package entity

import (
	"errors"
	"time"
)

var (
	ErrInvalidAccountID       = errors.New("invalid account id")
	ErrInvalidOperationTypeID = errors.New("invalid operation type id")
	ErrInvalidAmount          = errors.New("invalid amount")
	ErrInvalidAmountForOpType = errors.New("invalid amount for operation type")
	ErrInvalidEventDate       = errors.New("invalid event date")
)

type Transaction struct {
	ID              int       `json:"id"`
	AccountID       int       `json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"event_date"`
}

func (t *Transaction) Validate(opTypes OperationType) error {
	if t.AccountID <= 0 {
		return ErrInvalidAccountID
	}
	if t.Amount == 0 {
		return ErrInvalidAmount
	}
	if t.EventDate.IsZero() {
		return ErrInvalidEventDate
	}
	op, ok := opTypes[t.OperationTypeID]
	if !ok {
		return ErrInvalidOperationTypeID
	}
	if (op.PositiveAmount && t.Amount < 0) || (!op.PositiveAmount && t.Amount > 0) {
		return ErrInvalidAmountForOpType
	}
	return nil
}
