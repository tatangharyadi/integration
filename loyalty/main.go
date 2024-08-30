package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tatangharyadi/integration/loyalty/common/config"
	"github.com/tatangharyadi/integration/loyalty/pkg/voucherify"
)

func main() {
	env, logger := config.InitEnv()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	rs := voucherify.VoucherifyResource{
		Env:    env,
		Logger: logger,
	}
	r.Mount("/voucherify", rs.Routes())

	logger.Info().Msgf("Listening %s mode:%s", env.AppEnv, env.AppPort)
	addr := ":" + env.AppPort
	http.ListenAndServe(addr, r)
}
