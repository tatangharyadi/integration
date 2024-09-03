package models

type ParentProduct struct {
	ErpId       int     `json:"erp_id"`
	Sku         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Cost        float64 `json:"cost"`
}

type Product struct {
	ErpId       int           `json:"erp_id"`
	Sku         string        `json:"sku"`
	Barcode     string        `json:"barcode"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Cost        float64       `json:"cost"`
	Price       float64       `json:"price"`
	Parent      ParentProduct `json:"parent"`
}

type ProductId struct {
	ErpId      int    `json:"erp_id"`
	UpdateDate string `json:"write_date"`
}
