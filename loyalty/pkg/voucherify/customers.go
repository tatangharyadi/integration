package voucherify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tatangharyadi/integration/loyalty/model"
)

type VoucherifyMetadata struct {
	CardGuid                      string `json:"card_guiid"`
	EmployeeRedemptionDate        string `json:"employee_redemption_date"`
	EmployeeRedemptionPeriod      string `json:"employee_redemption_period"`
	EmployeeRedemptionMaxPeriod   int    `json:"employee_redemption_max_period"`
	EmployeeRedemptionTotalPeriod int    `json:"employee_redemption_total_period"`
}

type VoucherifyCustomer struct {
	SourceId string             `json:"source_id"`
	Name     string             `json:"name"`
	Phone    string             `json:"phone_number"`
	Metadata VoucherifyMetadata `json:"metadata"`
}

func MapCustomer(customer VoucherifyCustomer) model.Customer {
	return model.Customer{
		CustomerId: customer.SourceId,
		Name:       customer.Name,
		Phone:      customer.Phone,
		Credit: model.Credit{
			LastTransactionDate: customer.Metadata.EmployeeRedemptionDate,
			Period:              customer.Metadata.EmployeeRedemptionPeriod,
			Limit:               customer.Metadata.EmployeeRedemptionMaxPeriod,
			CurrentTotal:        customer.Metadata.EmployeeRedemptionTotalPeriod,
		},
	}
}

func (rs VoucherifyResource) GetCustomer(w http.ResponseWriter, r *http.Request) {
	customerId := chi.URLParam(r, "customerId")
	url := fmt.Sprintf("https://as1.api.voucherify.io/v1/customers/%s", customerId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("X-App-Id", rs.Env.VoucherifyId)
	req.Header.Set("X-App-Token", rs.Env.VoucherifySecretKey)

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
