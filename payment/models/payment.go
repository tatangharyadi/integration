package models

type PaymentMetadata struct {
	Token string `json:"token"`
}

type QrPayment struct {
	ReferenceId string          `json:"reference_id"`
	Type        string          `json:"type"`
	Currency    string          `json:"currency"`
	Amount      float64         `json:"amount"`
	QrString    string          `json:"qr_string"`
	Status      string          `json:"status"`
	Metadata    PaymentMetadata `json:"metadata"`
}
