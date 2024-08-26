package voucherify

import (
	"github.com/go-chi/chi/v5"
	"github.com/tatangharyadi/integration/loyalty/common/config"
)

type VoucherifyResource struct {
	Env *config.Env
}

func (rs VoucherifyResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/getcustomer/{customerId}", rs.GetCustomer)
	return r
}
