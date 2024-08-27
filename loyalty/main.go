package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	//"github.com/tatangharyadi/integration/loyalty/common/config"
	//"github.com/tatangharyadi/integration/loyalty/pkg/voucherify"
)

func main() {
	// env, log := config.InitEnv()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	// rs := voucherify.VoucherifyResource{
	// 	Env: env,
	// 	Log: log,
	// }
	// r.Mount("/voucherify", rs.Routes())

	addr := ":8080" //+ env.AppPort
	http.ListenAndServe(addr, r)
}
