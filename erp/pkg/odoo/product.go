package odoo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kolo/xmlrpc"
)

type paramGetProduct struct {
	FromDate string `json:"fromDate"`
	ToDate   string `json:"toDate"`
}

func (h Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var paramGetProduct paramGetProduct
	err := json.NewDecoder(r.Body).Decode(&paramGetProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	domainFilter := []any{
		[]any{"active", "=", true},
		[]any{"sale_ok", "=", true},
		[]any{"write_date", ">=", paramGetProduct.FromDate},
		[]any{"write_date", "<=", paramGetProduct.ToDate},
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
				"write_date",
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
