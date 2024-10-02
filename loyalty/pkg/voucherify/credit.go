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
)

type TransactionRequest struct {
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}

type VoucherifyMealBenefitMetadata struct {
	CreditBenefit VoucherifyCredit `json:"meal_benefit"`
}

type VoucherifyCreditBenefitMetadata struct {
	CreditBenefit VoucherifyCredit `json:"credit_benefit"`
}

type VoucherifyPersonalCreditMetadata struct {
	CreditBenefit VoucherifyCredit `json:"personal_credit"`
}

type VoucherifyMealBenefitTransaction struct {
	SourceId string                        `json:"source_id"`
	Metadata VoucherifyMealBenefitMetadata `json:"metadata"`
}

type VoucherifyCreditBenefitTransaction struct {
	SourceId string                          `json:"source_id"`
	Metadata VoucherifyCreditBenefitMetadata `json:"metadata"`
}

type VoucherifyPersonalCreditTransaction struct {
	SourceId string                           `json:"source_id"`
	Metadata VoucherifyPersonalCreditMetadata `json:"metadata"`
}

func getBalance(credit VoucherifyCredit) float64 {
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

	var VoucherifyCustomer VoucherifyCustomer
	if err := json.Unmarshal(body, &VoucherifyCustomer); err != nil {
		return models.Customer{}, err
	}

	customer := mapCustomer(VoucherifyCustomer)

	return customer, nil
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
	var payload []byte
	switch transaction.Type {
	case "MEALBENEFIT":
		if transaction.Amount > customer.MealBenefit.AvailableBalance {
			http.Error(w, "Insufficient balance", http.StatusBadRequest)
			return
		}

		var creditTransaction VoucherifyMealBenefitTransaction
		creditTransaction.SourceId = customer.Id
		creditTransaction.Metadata.CreditBenefit.Cycle = customer.MealBenefit.Cycle
		creditTransaction.Metadata.CreditBenefit.Limit = customer.MealBenefit.Limit
		creditTransaction.Metadata.CreditBenefit.Balance = customer.MealBenefit.Balance - transaction.Amount
		creditTransaction.Metadata.CreditBenefit.LastTransactionDate = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

		payload, err = json.Marshal(creditTransaction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "CREDITBENEFIT":
		if transaction.Amount > customer.CreditBenefit.AvailableBalance {
			http.Error(w, "Insufficient balance", http.StatusBadRequest)
			return
		}

		var creditTransaction VoucherifyCreditBenefitTransaction
		creditTransaction.SourceId = customer.Id
		creditTransaction.Metadata.CreditBenefit.Cycle = customer.CreditBenefit.Cycle
		creditTransaction.Metadata.CreditBenefit.Limit = customer.CreditBenefit.Limit
		creditTransaction.Metadata.CreditBenefit.Balance = customer.CreditBenefit.Balance - transaction.Amount
		creditTransaction.Metadata.CreditBenefit.LastTransactionDate = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

		payload, err = json.Marshal(creditTransaction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "PERSONALCREDIT":
		if transaction.Amount > customer.PersonalCredit.AvailableBalance {
			http.Error(w, "Insufficient balance", http.StatusBadRequest)
			return
		}

		var creditTransaction VoucherifyPersonalCreditTransaction
		creditTransaction.SourceId = customer.Id
		creditTransaction.Metadata.CreditBenefit.Cycle = customer.PersonalCredit.Cycle
		creditTransaction.Metadata.CreditBenefit.Limit = customer.PersonalCredit.Limit
		creditTransaction.Metadata.CreditBenefit.Balance = customer.PersonalCredit.Balance - transaction.Amount
		creditTransaction.Metadata.CreditBenefit.LastTransactionDate = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

		payload, err = json.Marshal(creditTransaction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

func DebitCustomer() {

}
