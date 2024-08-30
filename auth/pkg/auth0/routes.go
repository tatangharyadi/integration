package auth0

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/tatangharyadi/integration/auth/common/config"
)

type Auth0Resource struct {
	Env    *config.Env
	Logger zerolog.Logger
}

func (rs Auth0Resource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/gettoken", rs.GetToken)
	return r
}
