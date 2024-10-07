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
	Cycle               string  `json:"cycle,omitempty"`
	Limit               float64 `json:"limit,omitempty"`
	Balance             float64 `json:"balance,omitempty"`
	LastTransactionDate string  `json:"last_transaction_date,omitempty"`
}

type VoucherifyMetadata struct {
	EmployeeId     string            `json:"employee_id,omitempty"`
	MealBenefit    *VoucherifyCredit `json:"meal_benefit,omitempty"`
	CreditBenefit  *VoucherifyCredit `json:"credit_benefit,omitempty"`
	PersonalCredit *VoucherifyCredit `json:"personal_credit,omitempty"`
}

type VoucherifyCustomer struct {
	SourceId string             `json:"source_id"`
	Name     string             `json:"name,omitempty"`
	Email    string             `json:"email,omitempty"`
	Phone    string             `json:"phone,omitempty"`
	Metadata VoucherifyMetadata `json:"metadata"`
}

func mapCustomer(customer VoucherifyCustomer) models.Customer {
	return models.Customer{
		Id:         customer.SourceId,
		EmployeeId: customer.Metadata.EmployeeId,
		Name:       customer.Name,
		Email:      customer.Email,
		Phone:      customer.Phone,
		MealBenefit: models.Credit{
			Cycle:                customer.Metadata.MealBenefit.Cycle,
			Limit:                customer.Metadata.MealBenefit.Limit,
			Balance:              customer.Metadata.MealBenefit.Balance,
			TransactionTimestamp: customer.Metadata.MealBenefit.LastTransactionDate,
			AvailableBalance:     getBalance(*customer.Metadata.MealBenefit),
		},
		CreditBenefit: models.Credit{
			Cycle:                customer.Metadata.CreditBenefit.Cycle,
			Limit:                customer.Metadata.CreditBenefit.Limit,
			Balance:              customer.Metadata.CreditBenefit.Balance,
			TransactionTimestamp: customer.Metadata.CreditBenefit.LastTransactionDate,
			AvailableBalance:     getBalance(*customer.Metadata.CreditBenefit),
		},
		PersonalCredit: models.Credit{
			Cycle:                customer.Metadata.PersonalCredit.Cycle,
			Limit:                customer.Metadata.PersonalCredit.Limit,
			Balance:              customer.Metadata.PersonalCredit.Balance,
			TransactionTimestamp: customer.Metadata.PersonalCredit.LastTransactionDate,
			AvailableBalance:     getBalance(*customer.Metadata.PersonalCredit),
		},
	}
}

func mapVoucherify(customer models.Customer) VoucherifyCustomer {
	return VoucherifyCustomer{
		SourceId: customer.Id,
		Name:     customer.Name,
		Email:    customer.Email,
		Phone:    customer.Phone,
		Metadata: VoucherifyMetadata{
			EmployeeId: customer.EmployeeId,
			MealBenefit: &VoucherifyCredit{
				Cycle:               customer.MealBenefit.Cycle,
				Limit:               customer.MealBenefit.Limit,
				Balance:             customer.MealBenefit.Balance,
				LastTransactionDate: customer.MealBenefit.TransactionTimestamp,
			},
			CreditBenefit: &VoucherifyCredit{
				Cycle:               customer.CreditBenefit.Cycle,
				Limit:               customer.CreditBenefit.Limit,
				Balance:             customer.CreditBenefit.Balance,
				LastTransactionDate: customer.CreditBenefit.TransactionTimestamp,
			},
			PersonalCredit: &VoucherifyCredit{
				Cycle:               customer.PersonalCredit.Cycle,
				Limit:               customer.PersonalCredit.Limit,
				Balance:             customer.PersonalCredit.Balance,
				LastTransactionDate: customer.PersonalCredit.TransactionTimestamp,
			},
		},
	}
}

func getCustomer(h Handler, id string) (models.Customer, error) {
	url := fmt.Sprintf("%s/customers/%s", h.Env.LoyaltyUrl, id)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.Customer{}, err
	}

	req.Header.Set("X-App-Id", h.Env.VoucherifyId)
	req.Header.Set("X-App-Token", h.Env.VoucherifySecretKey)

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

func (h Handler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	customer, err := getCustomer(h, id)
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
	w.Write(resJson)
}

func importCustomer(h Handler, customer models.Customer,
	ch chan<- models.Customer, wg *sync.WaitGroup) interface{} {
	h.Logger.Info().Msgf("Importing customer %s", customer.Id)
	defer wg.Done()

	url := fmt.Sprintf("%s/customers", h.Env.LoyaltyUrl)
	client := &http.Client{}
	customer.MealBenefit.TransactionTimestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	customer.CreditBenefit.TransactionTimestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	customer.PersonalCredit.TransactionTimestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	voucherifyCustomer := mapVoucherify(customer)

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
	var VoucherifyCustomer VoucherifyCustomer
	if err := json.Unmarshal(body, &VoucherifyCustomer); err != nil {
		return err
	}

	customer = mapCustomer(VoucherifyCustomer)

	ch <- customer

	return nil
}

func (h Handler) ImportCustomers(w http.ResponseWriter, r *http.Request) {
	var customers models.Customers
	if err := json.NewDecoder(r.Body).Decode(&customers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ch := make(chan models.Customer)
	var wg sync.WaitGroup

	for _, customer := range customers.Customers {
		wg.Add(1)
		go importCustomer(h, customer, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	resCustomers := []models.Customer{}
	for result := range ch {

		resCustomers = append(resCustomers, result)
	}

	resJson, err := json.Marshal(resCustomers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}
