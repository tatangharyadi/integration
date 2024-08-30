package model

type Credit struct {
	LastTransactionDate string `json:"last_transaction_date"`
	Period              string `json:"period"`
	Limit               int    `json:"limit"`
	CurrentTotal        int    `json:"current_total"`
}

type Customer struct {
	CustomerId string `json:"customer_id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Credit     Credit `json:"credit"`
}
