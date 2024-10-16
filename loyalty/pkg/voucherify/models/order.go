package voucherify

type Order struct {
	SourceId            string  `json:"source_id"`
	Amount              float64 `json:"amount"`
	TotalDiscountAmount float64 `json:"total_discount_amount"`
}
