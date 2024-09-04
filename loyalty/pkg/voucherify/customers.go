package voucherify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tatangharyadi/integration/loyalty/models"
)

type VoucherifyCredit struct {
	Cycle               string  `json:"cycle"`
	Limit               float64 `json:"limit"`
	Balance             float64 `json:"balance"`
	LastTransactionDate string  `json:"last_transaction_date"`
}

type VoucherifyMetadata struct {
	CompanyBenefit VoucherifyCredit `json:"company_benefit"`
	PersonalCredit VoucherifyCredit `json:"personal_credit"`
}

type VoucherifyCustomer struct {
	SourceId string             `json:"source_id"`
	Name     string             `json:"name"`
	Email    string             `json:"email"`
	Phone    string             `json:"phone"`
	Metadata VoucherifyMetadata `json:"metadata"`
}

func MapCustomer(customer VoucherifyCustomer) models.Customer {
	return models.Customer{
		Id:    customer.SourceId,
		Name:  customer.Name,
		Email: customer.Email,
		Phone: customer.Phone,
		CompanyBenefit: models.Credit{
			Cycle:                customer.Metadata.CompanyBenefit.Cycle,
			Limit:                customer.Metadata.CompanyBenefit.Limit,
			Balance:              customer.Metadata.CompanyBenefit.Balance,
			TransactionTimestamp: customer.Metadata.CompanyBenefit.LastTransactionDate,
		},
		PersonalCredit: models.Credit{
			Cycle:                customer.Metadata.PersonalCredit.Cycle,
			Limit:                customer.Metadata.PersonalCredit.Limit,
			Balance:              customer.Metadata.PersonalCredit.Balance,
			TransactionTimestamp: customer.Metadata.PersonalCredit.LastTransactionDate,
		},
	}
}

func (h Handler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	url := fmt.Sprintf("%s/customers/%s", h.Env.LoyaltyUrl, id)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("X-App-Id", h.Env.VoucherifyId)
	req.Header.Set("X-App-Token", h.Env.VoucherifySecretKey)

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

	var VoucherifyCustomer VoucherifyCustomer
	if err := json.Unmarshal(body, &VoucherifyCustomer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	customer := MapCustomer(VoucherifyCustomer)
	resJson, err := json.Marshal(customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resJson)
}
