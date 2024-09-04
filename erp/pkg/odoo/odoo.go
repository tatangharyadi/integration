package odoo

import (
	"encoding/json"

	"github.com/kolo/xmlrpc"
)

func ReadAll(h Handler, uid int64, models *xmlrpc.Client, model string, method string,
	domain map[string]OdooField, domainFilter []any, domainField map[string]any) ([]byte, error) {
	var fields []string
	for _, value := range domain {
		fields = append(fields, value.name)
	}
	domainField["fields"] = fields

	var body []map[string]any
	if err := models.Call("execute_kw", []any{
		h.Env.OdooDb, uid, h.Env.OdooPassword,
		model, method,
		[]any{domainFilter},
		map[string]any{
			"fields": fields,
		},
	}, &body); err != nil {
		return []byte{}, err
	}

	for i, record := range body {
		for key, value := range record {
			if !domain[key].isBool && value == false {
				delete(body[i], key)
			}
		}
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return []byte{}, err
	}

	return jsonData, nil
}
