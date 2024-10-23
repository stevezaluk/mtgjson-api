package sealed_product

type SealedProductContents struct {
	Card     []SealedProductCard              `json:"card"`
	Deck     []SealedProductDeck              `json:"deck"`
	Other    []SealedProductOther             `json:"other"`
	Pack     []SealedProductPack              `json:"pack"`
	Sealed   []SealedProductSealed            `json:"sealed"`
	Variable map[string]SealedProductContents `json:"variable"`
}
