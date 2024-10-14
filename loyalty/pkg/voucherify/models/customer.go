package voucherify

type Credit struct {
	Cycle               string  `json:"cycle,omitempty"`
	Limit               float64 `json:"limit,omitempty"`
	Balance             float64 `json:"balance,omitempty"`
	LastTransactionDate string  `json:"last_transaction_date,omitempty"`
}

type Metadata struct {
	EmployeeId     string  `json:"employee_id,omitempty"`
	MealBenefit    *Credit `json:"meal_benefit,omitempty"`
	CreditBenefit  *Credit `json:"credit_benefit,omitempty"`
	PersonalCredit *Credit `json:"personal_credit,omitempty"`
}

type Customer struct {
	SourceId string   `json:"source_id"`
	Name     string   `json:"name,omitempty"`
	Email    string   `json:"email,omitempty"`
	Phone    string   `json:"phone,omitempty"`
	Metadata Metadata `json:"metadata"`
}
