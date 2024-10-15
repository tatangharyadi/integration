package voucherify

type Redeemables struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

type Session struct {
	Type string `json:"type"`
}

type Redeem struct {
	Redeemables []Redeemables `json:"redeemables"`
	Order       Order         `json:"order"`
	Customer    Customer      `json:"customer"`
	Session     Session       `json:"session"`
}

type Redemption struct {
	Order    Order    `json:"order"`
	Customer Customer `json:"customer"`
	Result   string   `json:"result"`
	Status   string   `json:"status"`
	Voucher  Voucher  `json:"voucher"`
}

type Redemptions struct {
	Redemptions []Redemption `json:"redemptions"`
}
