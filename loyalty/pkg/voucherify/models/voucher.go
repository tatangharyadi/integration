package voucherify

type VoucherDiscount struct {
	Type       string `json:"type"`
	PercentOff int    `json:"percent_off"`
}

type Voucher struct {
	Id              string           `json:"id,omitempty"`
	Code            string           `json:"code,omitempty"`
	Category        string           `json:"category,omitempty"`
	Type            string           `json:"type,omitempty"`
	Discount        *VoucherDiscount `json:"discount"`
	ValidationRules []string         `json:"validation_rules"`
	Active          bool             `json:"active,omitempty"`
}

type VoucherPublication struct {
	Voucher  string   `json:"voucher"`
	Customer Customer `json:"customer"`
}

type Vouchers struct {
	Vouchers []Voucher `json:"vouchers"`
}
