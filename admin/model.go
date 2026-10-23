package admin

type Product struct {
	ProductId string  `json:"product_id"`
	Name      string  `json:"name"`
	Catogary  string  `json:"catogary"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}
