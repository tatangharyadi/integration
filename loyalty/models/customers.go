package models

type Credit struct {
	Cycle                string  `json:"cycle"`
	Limit                float64 `json:"limit"`
	Balance              float64 `json:"balance"`
	TransactionTimestamp string  `json:"transaction_timestamp"`
	AvailableBalance     float64 `json:"available_balance"`
}

type Customer struct {
	Id             string `json:"id"`
	EmployeeId     string `json:"employee_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	MealBenefit    Credit `json:"meal_benefit"`
	CreditBenefit  Credit `json:"credit_benefit"`
	PersonalCredit Credit `json:"personal_credit"`
}

type Customers struct {
	Customers []Customer `json:"customers"`
}
