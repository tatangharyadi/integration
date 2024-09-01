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

	r.Get("/getcustomer/{customerId}", h.GetCustomer)
	return r
}
