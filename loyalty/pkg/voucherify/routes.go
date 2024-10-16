package voucherify

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/tatangharyadi/integration/loyalty/common/configs"
)

type Handler struct {
	Env    *configs.Env
	Logger zerolog.Logger
}

func (h Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/getcustomer/{id}", h.GetCustomer)
	r.Post("/importcustomers", h.ImportCustomers)
	r.Post("/redeem", h.Redeem)
	r.Post("/creditcustomer/{id}", h.CreditCustomer)
	r.Post("/debitcustomer/{id}", h.DebitCustomer)
	return r
}
