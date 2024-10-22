package sealed_product

import card "mtgjson/models/card/meta"

type SealedProduct struct {
	CardCount    int64                 `json:"cardCount"`
	Category     string                `json:"category"`
	Contents     SealedProductContents `json:"contents"`
	Identifiers  card.CardIdentifiers  `json:"identifiers"`
	Name         string                `json:"name"`
	ProductSize  int64                 `json:"productSize"`
	PurchaseUrls card.PurchaseUrls     `json:"purchaseUrls"`
	ReleaseDate  string                `json:"releaseDate"`
	Subtype      string                `json:"subtype"`
	UUID         string                `json:"uuid"`
}
