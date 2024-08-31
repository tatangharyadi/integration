package odoo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kolo/xmlrpc"
)

func (h Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	fromDate := r.URL.Query().Get("fromDate")
	toDate := r.URL.Query().Get("toDate")

	domainFilter := []any{
		[]any{"active", "=", true},
		[]any{"sale_ok", "=", true},
		[]any{"write_date", ">=", fromDate},
		[]any{"write_date", "<=", toDate},
	}
	if id != "*" {
		domainFilter = append(domainFilter, []any{"id", "=", id})
	}

	client, err := xmlrpc.NewClient(fmt.Sprintf("%s/xmlrpc/2/common", h.Env.ErpUrl), nil)
	if err != nil {
		log.Fatal(err)
	}

	var uid int64
	if err := client.Call("authenticate", []any{
		h.Env.OdooDb, h.Env.OdooUser, h.Env.OdooPassword,
		map[string]any{},
	}, &uid); err != nil {
		log.Fatal(err)
	}

	models, err := xmlrpc.NewClient(fmt.Sprintf("%s/xmlrpc/2/object", h.Env.ErpUrl), nil)
	if err != nil {
		log.Fatal(err)
	}
	var body []map[string]any
	if err := models.Call("execute_kw", []any{
		h.Env.OdooDb, uid, h.Env.OdooPassword,
		"product.product", "search_read",
		[]any{domainFilter},
		map[string]any{
			"fields": []string{
				"id",
				"barcode",
				"name",
				"description",
				"description_sale",
				"standard_price",
				"list_price",
			},
			"limit": 5,
		},
	}, &body); err != nil {
		log.Fatal(err)
	}

	resJson, err := json.Marshal(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}
