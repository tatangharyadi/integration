package auth0

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	config "github.com/tatangharyadi/integration/auth/common/configs"
)

type Handler struct {
	Env    *config.Env
	Logger zerolog.Logger
}

func (h Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/gettoken", h.GetToken)
	return r
}
