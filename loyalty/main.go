package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	config "github.com/tatangharyadi/integration/loyalty/common/configs"
	"github.com/tatangharyadi/integration/loyalty/pkg/voucherify"
)

func main() {
	env, logger := config.InitEnv()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World from loyalty"))
	})

	h := voucherify.Handler{
		Env:    env,
		Logger: logger,
	}
	r.Mount("/voucherify", h.Routes())

	logger.Info().Msgf("Listening %s mode:%s", env.AppEnv, env.AppPort)
	addr := ":" + env.AppPort
	http.ListenAndServe(addr, r)
}
