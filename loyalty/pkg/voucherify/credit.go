package voucherify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tatangharyadi/integration/loyalty/models"
	voucherify "github.com/tatangharyadi/integration/loyalty/pkg/voucherify/models"
)

type TransactionRequest struct {
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}

func getBalance(credit voucherify.Credit) float64 {
	today := time.Now().UTC()
	t, err := time.Parse("2006-01-02T15:04:05.000Z", credit.LastTransactionDate)
	if err != nil {
		return credit.Balance
	}

	if credit.Cycle == "DD" {
		if t.Before(today.AddDate(0, 0, -1)) {
			return credit.Limit
		}
	}
	return credit.Balance
}

func updateBalance(h Handler, payload []byte) (models.Customer, error) {
	url := fmt.Sprintf("%s/customers", h.Env.LoyaltyUrl)
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return models.Customer{}, err
	}

	req.Header.Set("X-App-Id", h.Env.VoucherifyId)
	req.Header.Set("X-App-Token", h.Env.VoucherifySecretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return models.Customer{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Customer{}, err
	}

	var VoucherifyCustomer voucherify.Customer
	if err := json.Unmarshal(body, &VoucherifyCustomer); err != nil {
		return models.Customer{}, err
	}

	customer := mapCustomer(VoucherifyCustomer)

	return customer, nil
}

func creditVoucherify(credit models.Credit, amount float64) (*voucherify.Credit, bool) {
	remainBalance := credit.AvailableBalance - amount
	if remainBalance < 0 {
		return &voucherify.Credit{}, false
	}
	var timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	return &voucherify.Credit{
		Cycle:               credit.Cycle,
		Limit:               credit.Limit,
		Balance:             remainBalance,
		LastTransactionDate: timestamp,
	}, true
}

func (h Handler) CreditCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var transaction TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	customer, err := getCustomer(h, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var creditTransaction voucherify.Customer
	var valid bool
	creditTransaction.SourceId = customer.Id
	switch transaction.Type {
	case "MEALBENEFIT":
		creditTransaction.Metadata.MealBenefit, valid = creditVoucherify(customer.MealBenefit, transaction.Amount)
		if !valid {
			http.Error(w, "Insufficient balance", http.StatusBadRequest)
			return
		}
	case "CREDITBENEFIT":
		creditTransaction.Metadata.CreditBenefit, valid = creditVoucherify(customer.CreditBenefit, transaction.Amount)
		if !valid {
			http.Error(w, "Insufficient balance", http.StatusBadRequest)
			return
		}
	case "PERSONALCREDIT":
		creditTransaction.Metadata.PersonalCredit, valid = creditVoucherify(customer.PersonalCredit, transaction.Amount)
		if !valid {
			http.Error(w, "Insufficient balance", http.StatusBadRequest)
			return
		}
	}

	var payload []byte
	payload, err = json.Marshal(creditTransaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	customer, err = updateBalance(h, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resJson, err := json.Marshal(customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}

func debitVoucherify(credit models.Credit, amount float64) *voucherify.Credit {
	remainBalance := credit.AvailableBalance + amount
	var timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	return &voucherify.Credit{
		Cycle:               credit.Cycle,
		Limit:               credit.Limit,
		Balance:             remainBalance,
		LastTransactionDate: timestamp,
	}
}

func (h Handler) DebitCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var transaction TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	customer, err := getCustomer(h, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var debitTransaction voucherify.Customer
	debitTransaction.SourceId = customer.Id
	switch transaction.Type {
	case "MEALBENEFIT":
		debitTransaction.Metadata.MealBenefit = debitVoucherify(customer.MealBenefit, transaction.Amount)
	case "CREDITBENEFIT":
		debitTransaction.Metadata.CreditBenefit = debitVoucherify(customer.CreditBenefit, transaction.Amount)
	case "PERSONALCREDIT":
		debitTransaction.Metadata.PersonalCredit = debitVoucherify(customer.PersonalCredit, transaction.Amount)
	}

	var payload []byte
	payload, err = json.Marshal(debitTransaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	customer, err = updateBalance(h, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resJson, err := json.Marshal(customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}
