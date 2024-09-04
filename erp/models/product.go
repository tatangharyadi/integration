package models

type ParentProduct struct {
	Id          int     `json:"id"`
	Sku         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Cost        float64 `json:"cost"`
}

type Product struct {
	Id          int           `json:"id"`
	Sku         string        `json:"sku"`
	Barcode     string        `json:"barcode"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Cost        float64       `json:"cost"`
	Price       float64       `json:"price"`
	Parent      ParentProduct `json:"parent"`
}

type ProductId struct {
	Id         int    `json:"id"`
	UpdateDate string `json:"write_date"`
}
