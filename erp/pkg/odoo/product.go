package odoo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kolo/xmlrpc"

	"github.com/tatangharyadi/integration/erp/models"
)

type ParamGetProduct struct {
	FromDate string `json:"fromDate"`
	ToDate   string `json:"toDate"`
}

type OdooField struct {
	name   string
	isBool bool
}

type OdooProduct struct {
	Id              int           `json:"id"`
	Template        []interface{} `json:"product_tmpl_id"`
	Barcode         string        `json:"barcode"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	DescriptionSale string        `json:"description_sale"`
	StandardPrice   float64       `json:"standard_price"`
	ListPrice       float64       `json:"list_price"`
	WriteDate       string        `json:"write_date"`
}

type OdooProductId struct {
	Id         int           `json:"id"`
	TemplateId []interface{} `json:"product_tmpl_id"`
	WriteDate  string        `json:"write_date"`
}

func MapProduct(odooProducts []OdooProduct, odooProductTemplates []OdooProduct) models.Product {
	return models.Product{
		Id:          odooProducts[0].Id,
		Sku:         odooProducts[0].Barcode,
		Barcode:     odooProducts[0].Barcode,
		Name:        odooProducts[0].Name,
		Description: odooProducts[0].Description,
		Cost:        odooProducts[0].StandardPrice,
		Price:       odooProducts[0].ListPrice,
		Parent: models.ParentProduct{
			Id:          odooProductTemplates[0].Id,
			Sku:         odooProductTemplates[0].Barcode,
			Name:        odooProductTemplates[0].Name,
			Description: odooProductTemplates[0].Description,
			Cost:        odooProductTemplates[0].StandardPrice,
		},
	}
}

func MapProductId(odooProductIds []OdooProductId) []models.ProductId {
	var productIds []models.ProductId
	for _, productId := range odooProductIds {
		productIds = append(productIds, models.ProductId{
			Id:              productId.Id,
			UpdateTimestamp: productId.WriteDate,
		})
	}
	return productIds
}

func (h Handler) GetProductIds(w http.ResponseWriter, r *http.Request) {
	var paramGetProduct ParamGetProduct
	err := json.NewDecoder(r.Body).Decode(&paramGetProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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

	domain := map[string]OdooField{
		"id":         {"id", false},
		"write_date": {"write_date", false},
	}
	domainFilter := []any{
		[]any{"active", "=", true},
		[]any{"sale_ok", "=", true},
		[]any{"write_date", ">=", paramGetProduct.FromDate},
		[]any{"write_date", "<=", paramGetProduct.ToDate},
	}
	domainFields := map[string]any{}
	jsonData, err := ReadAll(h, uid, models, "product.product", "search_read", domain, domainFilter, domainFields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var odooProductIds []OdooProductId
	if err := json.Unmarshal(jsonData, &odooProductIds); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	productIds := MapProductId(odooProductIds)
	resJson, err := json.Marshal(productIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}

func (h Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

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

	domain := map[string]OdooField{
		"id":               {"id", false},
		"product_tmpl_id":  {"product_tmpl_id", false},
		"barcode":          {"barcode", false},
		"name":             {"name", false},
		"description":      {"description", false},
		"description_sale": {"description_sale", false},
		"standard_price":   {"standard_price", false},
		"list_price":       {"list_price", false},
	}
	domainFilter := []any{
		[]any{"active", "=", true},
		[]any{"sale_ok", "=", true},
	}
	if id != ":id" {
		domainFilter = append(domainFilter, []any{"id", "=", id})
	}
	domainFields := map[string]any{
		"limit": 1,
	}

	jsonData, err := ReadAll(h, uid, models, "product.product", "search_read", domain, domainFilter, domainFields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var odooProducts []OdooProduct
	if err := json.Unmarshal(jsonData, &odooProducts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	domain = map[string]OdooField{
		"id":               {"id", false},
		"barcode":          {"barcode", false},
		"name":             {"name", false},
		"description":      {"description", false},
		"description_sale": {"description_sale", false},
		"standard_price":   {"standard_price", false},
	}
	domainFilter = []any{
		[]any{"active", "=", true},
		[]any{"sale_ok", "=", true},
		[]any{"id", "=", int(odooProducts[0].Template[0].(float64))},
	}
	domainFields = map[string]any{
		"limit": 1,
	}

	jsonData, err = ReadAll(h, uid, models, "product.template", "search_read", domain, domainFilter, domainFields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var odooProductTemplates []OdooProduct
	if err := json.Unmarshal(jsonData, &odooProductTemplates); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	products := MapProduct(odooProducts, odooProductTemplates)
	resJson, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}
