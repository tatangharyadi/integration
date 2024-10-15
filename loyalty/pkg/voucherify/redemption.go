package voucherify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tatangharyadi/integration/loyalty/models"
	voucherify "github.com/tatangharyadi/integration/loyalty/pkg/voucherify/models"
)

func mapVoucherifyRedeem(redemption models.Redeem) voucherify.Redeem {
	var redeemables []voucherify.Redeemables
	for _, redeemable := range redemption.Redeemables {
		redeemables = append(redeemables, voucherify.Redeemables{
			Object: redeemable.Object,
			Id:     redeemable.Id,
		})
	}
	return voucherify.Redeem{
		Redeemables: redeemables,
		Customer: voucherify.Customer{
			SourceId: redemption.Customer.Id,
		},
		Order: voucherify.Order{
			SourceId: redemption.Order.Id,
			Amount:   redemption.Order.Amount,
		},
		Session: voucherify.Session{
			Type: "LOCK",
		},
	}
}

func mapRedemption(voucherifyRedemptions []voucherify.Redemption) []models.Redemption {
	var redemptions []models.Redemption
	for _, voucherifyRedemption := range voucherifyRedemptions {
		redemptions = append(redemptions, models.Redemption{
			Voucher: models.Voucher{
				Code: voucherifyRedemption.Voucher.Code,
			},
			Order: models.Order{
				Id: voucherifyRedemption.Order.SourceId,
			},
			Status: voucherifyRedemption.Status,
		})
	}
	return redemptions
}

func (h Handler) Redeem(w http.ResponseWriter, r *http.Request) {
	var redemption models.Redeem
	if err := json.NewDecoder(r.Body).Decode(&redemption); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("%s/redemptions", h.Env.LoyaltyUrl)
	client := &http.Client{}

	voucherifyRedemption := mapVoucherifyRedeem(redemption)
	payload, err := json.Marshal(voucherifyRedemption)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("X-App-Id", h.Env.VoucherifyId)
	req.Header.Set("X-App-Token", h.Env.VoucherifySecretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	voucherifyRedemptions := voucherify.Redemptions{}
	if err := json.Unmarshal(body, &voucherifyRedemptions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redemptions := mapRedemption(voucherifyRedemptions.Redemptions)

	resJson, err := json.Marshal(redemptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}
