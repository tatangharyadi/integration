package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	config "github.com/tatangharyadi/integration/auth/common/configs"
	"github.com/tatangharyadi/integration/auth/pkg/auth0"
)

func main() {
	env, logger := config.InitEnv()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World from integration-auth"))
	})

	h := auth0.Handler{
		Env:    env,
		Logger: logger,
	}
	r.Mount("/auth0", h.Routes())

	logger.Info().Msgf("Listening %s mode:%s", env.AppEnv, env.AppPort)
	addr := ":" + env.AppPort
	http.ListenAndServe(addr, r)
}
