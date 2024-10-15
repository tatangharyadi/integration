package voucherify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tatangharyadi/integration/loyalty/models"
	voucherify "github.com/tatangharyadi/integration/loyalty/pkg/voucherify/models"
)

func mapCustomer(customer voucherify.Customer) models.Customer {
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

func mapVoucherify(customer models.Customer) voucherify.Customer {
	return voucherify.Customer{
		SourceId: customer.Id,
		Name:     customer.Name,
		Email:    customer.Email,
		Phone:    customer.Phone,
		Metadata: voucherify.Metadata{
			EmployeeId: customer.EmployeeId,
			MealBenefit: &voucherify.Credit{
				Cycle:               customer.MealBenefit.Cycle,
				Limit:               customer.MealBenefit.Limit,
				Balance:             customer.MealBenefit.Balance,
				LastTransactionDate: customer.MealBenefit.TransactionTimestamp,
			},
			CreditBenefit: &voucherify.Credit{
				Cycle:               customer.CreditBenefit.Cycle,
				Limit:               customer.CreditBenefit.Limit,
				Balance:             customer.CreditBenefit.Balance,
				LastTransactionDate: customer.CreditBenefit.TransactionTimestamp,
			},
			PersonalCredit: &voucherify.Credit{
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

func (h Handler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	customer, err := getCustomer(h, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	customer.Vouchers, _ = getVouchers(h, customer.Id)

	resJson, err := json.Marshal(customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}

func (h Handler) ImportCustomers(w http.ResponseWriter, r *http.Request) {
	var customers models.Customers
	if err := json.NewDecoder(r.Body).Decode(&customers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tasks := make([]Task, len(customers.Customers))
	for i, customer := range customers.Customers {
		tasks[i] = Task{
			handler:  h,
			customer: customer,
		}
	}

	wp := WorkerPool{
		Tasks:       tasks,
		concurrency: 5,
	}
	wp.Run()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
