package sealed_product

type SealedProductCard struct {
	Foil   bool   `json:"foil"`
	Name   string `json:"name"`
	Number string `json:"number"`
	Set    string `json:"set"`
	UUID   string `json:"uuid"`
}
