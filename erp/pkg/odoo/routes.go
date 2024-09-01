package odoo

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/tatangharyadi/integration/erp/common/configs"
)

type Handler struct {
	Env    *configs.Env
	Logger zerolog.Logger
}

func (h Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/getproduct/{id}", h.GetProduct)
	return r
}
