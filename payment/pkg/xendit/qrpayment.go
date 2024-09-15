package xendit

import (
	"encoding/json"
	"net/http"

	pubsub "github.com/tatangharyadi/pos-common/common/pubsub"
	message "github.com/tatangharyadi/pos-common/protobuf/message"
)

type XenditPaymentDetail struct {
	ReceiptId      string `json:"receipt_id"`
	Source         string `json:"source"`
	Name           string `json:"name"`
	AccountDetails string `json:"account_details"`
}

type XenditMetadata struct {
	Token string `json:"token"`
}

type XenditBasket struct {
	ReferenceId string  `json:"reference_id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Currency    string  `json:"currency"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Type        string  `json:"type"`
	Url         string  `json:"url"`
	Description string  `json:"description"`
	SubCategory string  `json:"sub_category"`
}

type XenditData struct {
	Id            string              `json:"id"`
	BusinessId    string              `json:"business_id"`
	Currency      string              `json:"currency"`
	Amount        float64             `json:"amount"`
	Status        string              `json:"status"`
	Created       string              `json:"created"`
	QrId          string              `json:"qr_id"`
	QrString      string              `json:"qr_string"`
	ReferenceId   string              `json:"reference_id"`
	Type          string              `json:"type"`
	ChannelCode   string              `json:"channel_code"`
	ExpireAt      string              `json:"expire_at"`
	Basket        []XenditBasket      `json:"basket"`
	Metadata      XenditMetadata      `json:"metadata"`
	PaymentDetail XenditPaymentDetail `json:"payment_detail"`
}

type XenditQrPayment struct {
	Event      string     `json:"event"`
	ApiVersion string     `json:"api_version"`
	BusinessId string     `json:"business_id"`
	Created    string     `json:"created"`
	Data       XenditData `json:"data"`
}

func MapQRPayment(xenditQrPayment XenditQrPayment) message.QrPayment {
	var data = message.QrPayment_Data{
		Id:          xenditQrPayment.Data.Id,
		BusinessId:  xenditQrPayment.Data.BusinessId,
		Currency:    xenditQrPayment.Data.Currency,
		Amount:      float32(xenditQrPayment.Data.Amount),
		Status:      xenditQrPayment.Data.Status,
		Created:     xenditQrPayment.Data.Created,
		QrId:        xenditQrPayment.Data.QrId,
		QrString:    xenditQrPayment.Data.QrString,
		ReferenceId: xenditQrPayment.Data.ReferenceId,
		Type:        xenditQrPayment.Data.Type,
		ChannelCode: xenditQrPayment.Data.ChannelCode,
		ExpiresAt:   xenditQrPayment.Data.ExpireAt,
		PaymentDetail: &message.QrPayment_Data_PaymentDetail{
			ReceiptId:      xenditQrPayment.Data.PaymentDetail.ReceiptId,
			Source:         xenditQrPayment.Data.PaymentDetail.Source,
			Name:           xenditQrPayment.Data.PaymentDetail.Name,
			AccountDetails: xenditQrPayment.Data.PaymentDetail.AccountDetails,
		},
	}

	return message.QrPayment{
		Token: xenditQrPayment.Data.Metadata.Token,
		Notification: &message.QrPayment_Notification{
			Title: "qr-payment",
			Body:  xenditQrPayment.Data.Status,
		},
		Data: &data,
	}
}

func (h Handler) CallbackQrPayment(w http.ResponseWriter, r *http.Request) {
	callbackToken := r.Header.Get("x-callback-token")
	if callbackToken != h.Env.XenditWebhookToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var xenditQrPayment XenditQrPayment
	err := json.NewDecoder(r.Body).Decode(&xenditQrPayment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	qrPayment := MapQRPayment(xenditQrPayment)
	if err := pubsub.PublishProtoMessages(h.Env.GCPProjectId, h.Env.QrPaymentTopic, &qrPayment); err != nil {
		h.Logger.Error().Err(err).Msg("Error publish message to Pub/Sub")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// w.Write(body)
}
