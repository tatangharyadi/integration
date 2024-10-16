package models

type Order struct {
	Id       string  `json:"id"`
	Amount   float64 `json:"amount"`
	Discount float64 `json:"discount"`
}
