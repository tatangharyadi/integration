package voucherify

import (
	"fmt"
	"io"
	"net/http"
)

func (rs VoucherifyResource) GetCustomer(w http.ResponseWriter, r *http.Request) {
	customerId := "cust_rImiBoXHZh1ivNX1dvbr5cqd"
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

	w.Write(body)
}
