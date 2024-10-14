package voucherify

type VoucherDiscount struct {
	DiscountType string `json:"type"`
	PercentOff   int    `json:"percent_off"`
}

type Voucher struct {
	Id              string          `json:"id,omitempty"`
	Code            string          `json:"code,omitempty"`
	Category        string          `json:"category"`
	Discount        VoucherDiscount `json:"discount"`
	VoucherType     string          `json:"type"`
	ValidationRules []string        `json:"validation_rules"`
}

type VoucherCustomer struct {
	SourceId string `json:"source_id"`
}

type VoucherPublication struct {
	Voucher  string          `json:"voucher"`
	Customer VoucherCustomer `json:"customer"`
}
