package entity

import "time"

type Transaction struct {
	ID              int
	AccountID       int
	OperationTypeID int
	Amount          float64
	EventDate       time.Time
}
