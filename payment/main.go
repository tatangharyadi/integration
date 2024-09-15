package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tatangharyadi/integration/payment/common/configs"
	"github.com/tatangharyadi/integration/payment/pkg/xendit"
)

func main() {
	env, logger := configs.InitEnv()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World from integration-payment"))
	})

	h := xendit.Handler{
		Env:    env,
		Logger: logger,
	}
	r.Mount("/xendit", h.Routes())

	logger.Info().Msgf("Listening %s mode:%s", env.AppEnv, env.AppPort)
	addr := ":" + env.AppPort
	http.ListenAndServe(addr, r)
}
