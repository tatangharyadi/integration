package models

type Credit struct {
	Cycle                string  `json:"cycle"`
	Limit                float64 `json:"limit"`
	Balance              float64 `json:"balance"`
	TransactionTimestamp string  `json:"transaction_timestamp"`
}

type Customer struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	CompanyBenefit Credit `json:"company_benefit"`
	PersonalCredit Credit `json:"personal_credit"`
}
