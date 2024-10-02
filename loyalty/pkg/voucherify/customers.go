package voucherify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

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
	EmployeeId     string           `json:"employee_id"`
	MealBenefit    VoucherifyCredit `json:"meal_benefit"`
	CreditBenefit  VoucherifyCredit `json:"credit_benefit"`
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
		Id:         customer.SourceId,
		EmployeeId: customer.Metadata.EmployeeId,
		Name:       customer.Name,
		Email:      customer.Email,
		Phone:      customer.Phone,
		MealBenefit: models.Credit{
			Cycle:                customer.Metadata.CreditBenefit.Cycle,
			Limit:                customer.Metadata.CreditBenefit.Limit,
			Balance:              customer.Metadata.CreditBenefit.Balance,
			TransactionTimestamp: customer.Metadata.CreditBenefit.LastTransactionDate,
		},
		CreditBenefit: models.Credit{
			Cycle:                customer.Metadata.CreditBenefit.Cycle,
			Limit:                customer.Metadata.CreditBenefit.Limit,
			Balance:              customer.Metadata.CreditBenefit.Balance,
			TransactionTimestamp: customer.Metadata.CreditBenefit.LastTransactionDate,
		},
		PersonalCredit: models.Credit{
			Cycle:                customer.Metadata.PersonalCredit.Cycle,
			Limit:                customer.Metadata.PersonalCredit.Limit,
			Balance:              customer.Metadata.PersonalCredit.Balance,
			TransactionTimestamp: customer.Metadata.PersonalCredit.LastTransactionDate,
		},
	}
}

func MapVoucherify(customer models.Customer) VoucherifyCustomer {
	return VoucherifyCustomer{
		SourceId: customer.Id,
		Name:     customer.Name,
		Email:    customer.Email,
		Phone:    customer.Phone,
		Metadata: VoucherifyMetadata{
			EmployeeId: customer.EmployeeId,
			MealBenefit: VoucherifyCredit{
				Cycle:               customer.MealBenefit.Cycle,
				Limit:               customer.MealBenefit.Limit,
				Balance:             customer.MealBenefit.Balance,
				LastTransactionDate: time.Now().Format("2006-01-02T15:04:05.000Z"),
			},
			CreditBenefit: VoucherifyCredit{
				Cycle:               customer.CreditBenefit.Cycle,
				Limit:               customer.CreditBenefit.Limit,
				Balance:             customer.CreditBenefit.Balance,
				LastTransactionDate: time.Now().Format("2006-01-02T15:04:05.000Z"),
			},
			PersonalCredit: VoucherifyCredit{
				Cycle:               customer.PersonalCredit.Cycle,
				Limit:               customer.PersonalCredit.Limit,
				Balance:             customer.PersonalCredit.Balance,
				LastTransactionDate: time.Now().Format("2006-01-02T15:04:05.000Z"),
			},
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

func ImportCustomer(h Handler, customer models.Customer,
	ch chan<- []byte, wg *sync.WaitGroup) interface{} {
	h.Logger.Info().Msgf("Importing customer %s", customer.Id)
	defer wg.Done()

	url := fmt.Sprintf("%s/customers", h.Env.LoyaltyUrl)
	client := &http.Client{}
	voucherifyCustomer := MapVoucherify(customer)
	customerBytes, err := json.Marshal(voucherifyCustomer)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(customerBytes))
	if err != nil {
		return err
	}

	req.Header.Set("X-App-Id", h.Env.VoucherifyId)
	req.Header.Set("X-App-Token", h.Env.VoucherifySecretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ch <- body

	return nil
}

func (h Handler) ImportCustomers(w http.ResponseWriter, r *http.Request) {
	var customers models.Customers
	if err := json.NewDecoder(r.Body).Decode(&customers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ch := make(chan []byte)
	var wg sync.WaitGroup

	for _, customer := range customers.Customers {
		wg.Add(1)
		go ImportCustomer(h, customer, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var vaucherifyCustomers []VoucherifyCustomer
	for result := range ch {
		var voucherifyCustomer VoucherifyCustomer
		if err := json.Unmarshal(result, &voucherifyCustomer); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		vaucherifyCustomers = append(vaucherifyCustomers, voucherifyCustomer)
	}

	resJson, err := json.Marshal(vaucherifyCustomers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}
