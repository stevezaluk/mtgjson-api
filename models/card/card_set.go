package card

import (
	"go.mongodb.org/mongo-driver/bson"
	"mtgjson/context"
	"mtgjson/errors"
	"mtgjson/models/card/meta"
	"regexp"
)

type CardSet struct {
	AsciiName               string                `json:"asciiName"`
	AttractionLights        []string              `json:"attractionLights"`
	ColorIdentity           []string              `json:"colorIdentity"`
	ColorIndicator          []string              `json:"colorIndicator"`
	Colors                  []string              `json:"colors"`
	ConvertedManaCost       int                   `json:"convertedManaCost"`
	Defense                 string                `json:"defense"`
	EDHRecRank              int                   `json:"edhrecRank"`
	EDHRecSaltiness         float64               `json:"edhrecSaltiness"`
	FaceConvertedManaCost   int                   `json:"faceConvertedManaCost"`
	FaceManaValue           int                   `json:"faceManaValue"`
	FaceName                string                `json:"faceName"`
	FirstPrinting           string                `json:"firstPrinting"`
	ForeignData             []card.ForeignData    `json:"foreignData"`
	Hand                    string                `json:"hand"`
	HasAlternativeDeckLimit bool                  `json:"hasAlternativeDeckLimit"`
	Identifiers             card.CardIdentifiers  `json:"identifiers"`
	IsFunny                 bool                  `json:"isFunny"`
	IsReserved              bool                  `json:"isReserved"`
	Keywords                []string              `json:"keywords"`
	Layout                  string                `json:"layout"`
	LeadershipSkills        card.LeadershipSkills `json:"leadershipSkills"`
	Legalities              card.CardLegalities   `json:"legalities"`
	Life                    string                `json:"life"`
	Loyalty                 string                `json:"loyalty"`
	ManaCost                string                `json:"manaCost"`
	ManaValue               int                   `json:"manaValue"`
	Name                    string                `json:"name"`
	Power                   string                `json:"power"`
	Printings               []string              `json:"printings"`
	PurchaseUrls            card.PurchaseUrls     `json:"purchaseUrls"`
	RelatedCards            card.RelatedCards     `json:"relatedCards"`
	Rulings                 card.CardRulings      `json:"rulings"`
	Side                    string                `json:"side"`
	Subsets                 []string              `json:"subsets"`
	Subtypes                []string              `json:"subtypes"`
	Supertypes              []string              `json:"supertypes"`
	Text                    string                `json:"text"`
	Toughness               string                `json:"toughness"`
	Type                    string                `json:"type"`
	Types                   []string              `json:"types"`
	Artist                  string                `json:"artist"`
	ArtistIds               []string              `json:"artistIds"`
	Availability            []string              `json:"availability"`
	BoosterTypes            []string              `json:"boosterTypes"`
	BorderColor             string                `json:"borderColor"`
	CardParts               []string              `json:"cardParts"`
	DuelDeck                string                `json:"duelDeck"`
	FaceFlavorName          string                `json:"faceFlavorName"`
	Finishes                []string              `json:"finishes"`
	FlavorName              string                `json:"flavorName"`
	FrameEffects            []string              `json:"frameEffects"`
	FrameVersion            string                `json:"frameVersion"`
	HasContentWarning       bool                  `json:"hasContentWarning"`
	HasFoil                 bool                  `json:"hasFoil"`
	HasNonFoil              bool                  `json:"hasNonFoil"`
	IsAlternative           bool                  `json:"isAlternative"`
	IsFullArt               bool                  `json:"isFullArt"`
	IsOnlineOnly            bool                  `json:"isOnlineOnly"`
	IsOversized             bool                  `json:"isOversized"`
	IsPromo                 bool                  `json:"isPromo"`
	IsRebalanced            bool                  `json:"isRebalanced"`
	IsReprint               bool                  `json:"isReprint"`
	IsStarter               bool                  `json:"isStarter"`
	IsStorySpotlight        bool                  `json:"isStorySpotlight"`
	IsTextless              bool                  `json:"isTextless"`
	IsTimeshifted           bool                  `json:"isTimeshifted"`
	Language                string                `json:"language"`
	Number                  string                `json:"number"`
	OriginalPrintings       []string              `json:"originalPrintings"`
	OriginalReleaseDate     string                `json:"originalReleaseDate"`
	OriginalText            string                `json:"originalText"`
	OriginalType            string                `json:"originalType"`
	OtherFaceIds            []string              `json:"otherFaceIds"`
	PromoTypes              []string              `json:"promoTypes"`
	Rarity                  string                `json:"rarity"`
	RebalancedPrintings     []string              `json:"rebalancedPrintings"`
	SecurityStamp           string                `json:"securityStamp"`
	SetCode                 string                `json:"setCode"`
	Signature               string                `json:"signature"`
	SourceProducts          card.SourceProducts   `json:"sourceProducts"`
	UUID                    string                `json:"uuid"`
	Variations              []string              `json:"variations"`
	Watermark               []string              `json:"watermark"`
}

func ValidateUUID(uuid string) bool {
	var ret = false
	var pattern = `^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`

	re := regexp.MustCompile(pattern)
	if re.MatchString(uuid) {
		ret = true
	}

	return ret
}

func IterCards(cards []string) []CardSet {
	var ret []CardSet
	for i := 0; i < len(cards); i++ {
		uuid := cards[i]

		card, err := GetCard(uuid)
		if err != nil {
			continue
		}

		ret = append(ret, card)
	}

	return ret
}

func GetCard(uuid string) (CardSet, error) {
	var result CardSet

	if !ValidateUUID(uuid) {
		return result, errors.ErrInvalidUUID
	}

	var database = context.GetDatabase()

	query := bson.M{"identifiers.mtgjsonV4Id": uuid}
	results := database.Find("card", query, &result)
	if results == nil {
		return result, errors.ErrNoCard
	}

	return result, nil
}

func IndexCards(limit int64) ([]CardSet, error) {
	var result []CardSet

	var database = context.GetDatabase()

	results := database.Index("card", limit, &result)
	if results == nil {
		return result, errors.ErrNoCards
	}

	return result, nil

}
