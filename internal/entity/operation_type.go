package entity

type OperationType map[int]Operation

type Operation struct {
	Description    string `json:"description"`
	PositiveAmount bool   `json:"positive_amount"`
}
