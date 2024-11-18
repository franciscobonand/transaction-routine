package entity

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	ErrMissingDocumentNumber = errors.New("missing document number")
)

type Account struct {
	ID                   int             `json:"id"`
	DocumentNumber       string          `json:"document_number"`
	AvailableCreditLimit decimal.Decimal `json:"available_credit_limit"`
}

type AccountFilter struct {
	ID             *int    `json:"id"`
	DocumentNumber *string `json:"document_number"`
}

func (a Account) ToFilter() AccountFilter {
	filter := AccountFilter{}
	if a.ID != 0 {
		filter.ID = &a.ID
	}
	if a.DocumentNumber != "" {
		filter.DocumentNumber = &a.DocumentNumber
	}
	return filter
}
