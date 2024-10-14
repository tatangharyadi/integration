package voucherify

type VoucherDiscount struct {
	Type       string `json:"type"`
	PercentOff int    `json:"percent_off"`
}

type Voucher struct {
	Id              string          `json:"id,omitempty"`
	Code            string          `json:"code,omitempty"`
	Category        string          `json:"category"`
	Type            string          `json:"type"`
	Discount        VoucherDiscount `json:"discount"`
	ValidationRules []string        `json:"validation_rules"`
	Active          bool            `json:"active"`
}

type VoucherCustomer struct {
	SourceId string `json:"source_id"`
}

type VoucherPublication struct {
	Voucher  string          `json:"voucher"`
	Customer VoucherCustomer `json:"customer"`
}

type Vouchers struct {
	Vouchers []Voucher `json:"vouchers"`
}
