package voucherify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	netUrl "net/url"

	"github.com/tatangharyadi/integration/loyalty/models"
	voucherify "github.com/tatangharyadi/integration/loyalty/pkg/voucherify/models"
)

func mapVoucher(voucher voucherify.Voucher) models.Voucher {
	return models.Voucher{
		Code:     voucher.Code,
		Category: voucher.Category,
		Type:     voucher.Type,
		Discount: models.VoucherDiscount{
			Type:       voucher.Discount.Type,
			PercentOff: voucher.Discount.PercentOff,
		},
		Active: voucher.Active,
	}
}

func getVouchers(h Handler, id string) ([]models.Voucher, error) {
	url := fmt.Sprintf("%s/vouchers", h.Env.LoyaltyUrl)
	paramUrl, err := netUrl.Parse(url)
	if err != nil {
		return []models.Voucher{}, err
	}

	params := netUrl.Values{}
	params.Add("customer", id)
	paramUrl.RawQuery = params.Encode()

	client := &http.Client{}
	req, err := http.NewRequest("GET", paramUrl.String(), nil)
	if err != nil {
		return []models.Voucher{}, err
	}

	req.Header.Set("X-App-Id", h.Env.VoucherifyId)
	req.Header.Set("X-App-Token", h.Env.VoucherifySecretKey)

	resp, err := client.Do(req)
	if err != nil {
		return []models.Voucher{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []models.Voucher{}, err
	}

	var voucherifyVoucher voucherify.Vouchers
	if err := json.Unmarshal(body, &voucherifyVoucher); err != nil {
		return []models.Voucher{}, err
	}

	var vouchers []models.Voucher
	for _, voucher := range voucherifyVoucher.Vouchers {
		vouchers = append(vouchers, mapVoucher(voucher))
	}

	return vouchers, nil
}
