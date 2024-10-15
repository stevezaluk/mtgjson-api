package models

import (
	"mtgjson/models/card/meta"
)

type CardSet struct {
	CardAtomic
	Artist              string                `json:"artist"`
	ArtistIds           []string              `json:"artistIds"`
	Availability        []string              `json:"availability"`
	BoosterTypes        []string              `json:"boosterTypes"`
	BorderColor         string                `json:"borderColor"`
	CardParts           []string              `json:"cardParts"`
	DuelDeck            string                `json:"duelDeck"`
	FaceFlavorName      string                `json:"faceFlavorName"`
	Finishes            []string              `json:"finishes"`
	FlavorName          string                `json:"flavorName"`
	FrameEffects        []string              `json:"FrameEffects"`
	FrameVersion        string                `json:"frameVersion"`
	HasContentWarning   bool                  `json:"hasContentWarning"`
	HasFoil             bool                  `json:"hasFoil"`
	HasNonFoil          bool                  `json:"hasNonFoil"`
	IsAlternative       bool                  `json:"isAlternative"`
	IsFullArt           bool                  `json:"isFullArt"`
	IsOnlineOnly        bool                  `json:"isOnlineOnly"`
	IsOversized         bool                  `json:"isOversized"`
	IsPromo             bool                  `json:"isPromo"`
	IsRebalanced        bool                  `json:"isRebalanced"`
	IsReprint           bool                  `json:"isReprint"`
	IsStarter           bool                  `json:"isStarter"`
	IsStorySpotlight    bool                  `json:"isStorySpotlight"`
	IsTextless          bool                  `json:"isTextless"`
	IsTimeshifted       bool                  `json:"isTimeshifted"`
	Language            string                `json:"language"`
	Number              string                `json:"number"`
	OriginalPrintings   []string              `json:"originalPrintings"`
	OriginalReleaseDate string                `json:"originalReleaseDate"`
	OriginalText        string                `json:"originalText"`
	OriginalType        string                `json:"originalType"`
	OtherFaceIds        []string              `json:"otherFaceIds"`
	PromoTypes          []string              `json:"promoTypes"`
	Rarity              string                `json:"rarity"`
	RebalancedPrintings []string              `json:"rebalancedPrintings"`
	SecurityStamp       string                `json:"securityStamp"`
	SetCode             string                `json:"setCode"`
	Signature           string                `json:"signature"`
	SourceProducts      models.SourceProducts `json:"sourceProducts"`
	UUID                string                `json:"uuid"`
	Variations          []string              `json:"variations"`
	Watermark           []string              `json:"watermark"`
}
