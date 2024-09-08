package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tatangharyadi/integration/database/common/configs"
	"github.com/tatangharyadi/integration/database/pkg/yummyos"
)

func main() {
	env, logger := configs.InitEnv()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World from integration-database"))
	})

	h := yummyos.Handler{
		Env:    env,
		Logger: logger,
	}
	r.Mount("/yummyos", h.Routes())

	logger.Info().Msgf("Listening %s mode:%s", env.AppEnv, env.AppPort)
	addr := ":" + env.AppPort
	http.ListenAndServe(addr, r)
}
