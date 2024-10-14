package models

type VoucherDiscount struct {
	Type       string `json:"type"`
	PercentOff int    `json:"percent_off"`
}

type Voucher struct {
	Code     string          `json:"code"`
	Category string          `json:"category"`
	Type     string          `json:"type"`
	Discount VoucherDiscount `json:"discount"`
	Active   bool            `json:"active"`
}

type Credit struct {
	Cycle                string  `json:"cycle"`
	Limit                float64 `json:"limit"`
	Balance              float64 `json:"balance"`
	TransactionTimestamp string  `json:"transaction_timestamp"`
	AvailableBalance     float64 `json:"available_balance"`
}

type Customer struct {
	Id             string    `json:"id"`
	EmployeeId     string    `json:"employee_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	MealBenefit    Credit    `json:"meal_benefit"`
	CreditBenefit  Credit    `json:"credit_benefit"`
	PersonalCredit Credit    `json:"personal_credit"`
	Vouchers       []Voucher `json:"vouchers"`
}
