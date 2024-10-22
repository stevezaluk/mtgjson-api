package set

import (
	"mtgjson/models/booster"
	"mtgjson/models/sealed_product"
	"mtgjson/models/set/meta"
)

type Set struct {
	BaseSetSize      int64                            `json:"baseSetSize"`
	Block            string                           `json:"block"`
	Booster          map[string]booster.BoosterConfig `json:"booster"`
	Cards            []string                         `json:"cards"`
	CardsphereId     int64                            `json:"cardsphereId"`
	Code             string                           `json:"code"`
	CodeV3           string                           `json:"codeV3"`
	Decks            []string                         `json:"decks"`
	IsForeignOnly    bool                             `json:"isForeignOnly"`
	IsFoilOnly       bool                             `json:"isFoilOnly"`
	IsNonFoilOnly    bool                             `json:"isNonFoilOnly"`
	IsOnlineOnly     bool                             `json:"isOnlineOnly"`
	IsPaperOnly      bool                             `json:"isPaperOnly"`
	IsPartialPreview bool                             `json:"isPartialPreview"`
	KeyruneCode      string                           `json:"keyruneCode"`
	Languages        []string                         `json:"languages"`
	McmId            int64                            `json:"mcmId"`
	McmIdExtras      int64                            `json:"mcmIdExtras"`
	McmName          string                           `json:"mcmName"`
	MTGOCode         string                           `json:"mtgoCode"`
	Name             string                           `json:"name"`
	ParentCode       string                           `json:"parentCode"`
	ReleaseDate      string                           `json:"releaseDate"`
	SealedProduct    []sealed_product.SealedProduct   `json:"sealedProduct"`
	TCGPlayerGroupId int64                            `json:"tcgPlayerGroupId"`
	Tokens           []string                         `json:"tokens"`
	TokenSetCode     string                           `json:"tokenSetCode"`
	TokenSetSize     int64                            `json:"tokenSetSize"`
	Translations     set.Translations                 `json:"translations"`
	Type             string                           `json:"type"`
}
