package xendit

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/tatangharyadi/integration/payment/common/configs"
)

type Handler struct {
	Env    *configs.Env
	Logger zerolog.Logger
}

func (h Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/createqrpayment", h.CreateQrPayment)
	r.Post("/callbackqrpayment", h.CallbackQrPayment)
	return r
}
