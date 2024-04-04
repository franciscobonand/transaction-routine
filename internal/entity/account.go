package entity

type Account struct {
	ID             int    `json:"id"`
	DocumentNumber string `json:"document_number"`
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
