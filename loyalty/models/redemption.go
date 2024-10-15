package models

type Redeemables struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

type Redeem struct {
	Redeemables []Redeemables `json:"redeemables"`
	Order       Order         `json:"order"`
	Customer    Customer      `json:"customer"`
}

type Redemption struct {
	Voucher Voucher `json:"voucher"`
	Order   Order   `json:"order"`
	Status  string  `json:"status"`
}

type Redemptions struct {
	Redemptions []Redemption `json:"redemptions"`
}
