package models

type Credit struct {
	Cycle               string  `json:"cycle"`
	Limit               float64 `json:"limit"`
	Balance             float64 `json:"balance"`
	LastTransactionDate string  `json:"last_transaction_date"`
}

type Customer struct {
	CustomerId     string `json:"customer_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	CompanyBenefit Credit `json:"company_benefit"`
	PersonalCredit Credit `json:"personal_credit"`
}
