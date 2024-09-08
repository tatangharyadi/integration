package yummyos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (h Handler) GetPlaceProducts(w http.ResponseWriter, r *http.Request) {
	query := `
		WITH 
		cteModifier AS
		(SELECT
		products.id,
		sku,
		name,
		jsonb_agg(
			jsonb_build_object(
			'currency_code', currency_code,
			'price', price,
			'price_effective_time', COALESCE(price_effective_time, ''),
			'price_expire_time', COALESCE(price_expire_time, '')
			)
		) as price_infos
		FROM products
		LEFT JOIN price_infos ON
		price_infos.product_id = products.id
		WHERE
		TYPE = 'MODIFIER'
		GROUP BY
		products.id),
		cteModifierCollection AS 
		(SELECT
		mc.id AS id,
		mc.name,
		selection_min,
		selection_max,
		jsonb_agg(
			jsonb_build_object(
			'sku', cteModifier.sku,
			'name', cteModifier.name,
			'price_infos', price_infos
			)
		) AS modifier
		FROM products mc
		LEFT JOIN product_relations ON
		product_relations.parent_product_id = mc.id
		LEFT JOIN cteModifier ON
		cteModifier.id = product_relations.product_id
		WHERE
		mc.type = 'MODIFIER_COLLECTION'
		GROUP BY
		mc.id, mc.name, selection_min, selection_max),
		cteProduct AS
		(SELECT
		p.id as id,
		p.sku,
		p.gtin as barcode,
		p.name,
		description,
		cost,
		jsonb_agg(
			jsonb_build_object(
			'name', cteModifierCollection.name,
			'selection_min', cteModifierCollection.selection_min,
			'selection_max', cteModifierCollection.selection_max,
			'modifier', cteModifierCollection.modifier
			)
		) as modifier_collection
		FROM products p
		LEFT JOIN product_relations ON
		product_relations.parent_product_id = p.id
		LEFT JOIN cteModifierCollection ON
		cteModifierCollection.id = product_relations.product_id
		WHERE
		p.type = 'PRODUCT_LOCAL'
		GROUP BY
		p.id),
		cteLocalPrice AS
		(SELECT
		local_inventories.id,
		jsonb_agg(
			jsonb_build_object(
			'currency_code', currency_code,
			'price', price,
			'price_effective_time', COALESCE(price_effective_time, ''),
			'price_expire_time', COALESCE(price_expire_time, '')
			)
		) as price_infos
		FROM local_inventories
		INNER JOIN price_infos ON
		price_infos.id = local_inventories.price_info_id
		GROUP BY
		local_inventories.id),
		cteLocal AS
		(SELECT
		place_id,
		local_inventories.product_id AS id,
		sku,
		barcode,
		COALESCE(NULLIF(local_inventories.name, ''), cteProduct.name) AS name,
		description,
		image,
		cost,
		price_infos,
		modifier_collection,
		local_inventories.availability,
		local_inventories.updated_at AS update_timestamp
		FROM local_inventories
		INNER JOIN cteProduct ON
		cteProduct.id = local_inventories.product_id
		INNER JOIN cteLocalPrice ON
		cteLocalPrice.id = local_inventories.id
		WHERE
		local_inventories.deleted_at IS NULL
		AND sales_channel LIKE '%POS%')
		SELECT
		row_to_json(t)
		FROM (
		SELECT
			cteLocal.id,
			sku,
			barcode,
			cteLocal.name,
			cteLocal.description,
			cteLocal.image,
			cost,
			price_infos,
			modifier_collection,
			availability,
			update_timestamp
		FROM cteLocal
		INNER JOIN places ON
			places.id = cteLocal.place_id
		WHERE
			cteLocal.id = $1
		) t;
	`

	id := chi.URLParam(r, "id")
	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s database=%s",
		h.Env.DbInstanceHost, h.Env.DbInstancePort,
		h.Env.DbUser, h.Env.DbPassword, h.Env.DbName)

	config, err := pgxpool.ParseConfig(dbURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute

	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := dbPool.Query(context.Background(), query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []json.RawMessage
	for rows.Next() {
		var jsonData json.RawMessage
		if err := rows.Scan(&jsonData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, jsonData)
	}

	if rows.Err() != nil {
		http.Error(w, rows.Err().Error(), http.StatusInternalServerError)
		return
	}

	resJson, err := json.Marshal(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}
