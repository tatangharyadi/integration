package voucherify

type Order struct {
	SourceId string  `json:"source_id"`
	Amount   float64 `json:"amount,omitempty"`
}
