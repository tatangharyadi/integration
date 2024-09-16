package models

type QrPayment struct {
	Token       string  `json:"token"`
	ReferenceId string  `json:"reference_id"`
	Type        string  `json:"type"`
	Currency    string  `json:"currency"`
	Amount      float64 `json:"amount"`
	QrString    string  `json:"qr_string"`
	Status      string  `json:"status"`
}
