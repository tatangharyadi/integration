package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tatangharyadi/integration/auth/common/config"
	"github.com/tatangharyadi/integration/auth/pkg/auth0"
)

func main() {
	env, logger := config.InitEnv()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World from auth"))
	})

	rs := auth0.Auth0Resource{
		Env:    env,
		Logger: logger,
	}
	r.Mount("/auth0", rs.Routes())

	logger.Info().Msgf("Listening %s mode:%s", env.AppEnv, env.AppPort)
	addr := ":" + env.AppPort
	http.ListenAndServe(addr, r)
}
