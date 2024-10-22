package card

import (
	"go.mongodb.org/mongo-driver/bson"
	"mtgjson/context"
	"mtgjson/errors"
	"mtgjson/models/card/meta"
	"regexp"
)

/*
Card - A model representing an MTGJSON Card

Ommiting card descriptions for brevity.
See: https://mtgjson.com/data-models/card/card-set/
*/
type Card struct {
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

/*
ValidateUUID - Ensure that the passed UUID is valid

Paremeters:
uuid (string) - The UUID you want to validate

Returns:
ret (bool) - True if the UUID is valid, false if it is not
*/
func ValidateUUID(uuid string) bool {
	var ret = false
	var pattern = `^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`

	re := regexp.MustCompile(pattern)
	if re.MatchString(uuid) {
		ret = true
	}

	return ret
}

/*
ValidateCards - Ensure a list of cards both exist and are valid UUID's

Paremeters:
uuids (array[string]) - A list of mtgjsonV4 UUID's to validate

Returns:
result (bool) - True if all cards passed validation, False if they didnt
invalidCards (array[string]) - Values that are not valid UUID's
noExistCards (array[string]) - Cards that do not exist in Mongo
*/
func ValidateCards(uuids []string) (bool, []string, []string) {
	var invalidCards []string // cards that failed UUID validation
	var noExistCards []string // cards that do not exist in Mongo
	var result = true

	for _, uuid := range uuids {
		_, err := GetCard(uuid)
		if err == errors.ErrNoCard {
			result = false
			noExistCards = append(noExistCards, uuid)
		} else if err == errors.ErrInvalidUUID {
			result = false
			invalidCards = append(invalidCards, uuid)
		}
	}

	return result, invalidCards, noExistCards
}

/*
GetCards - Takes a list of UUID's and returns card models for them

Paramters:
cards (slice[string]) - A list of UUID's you want card models for

Returns
ret (slice[card.Card]) - A list of card models
*/
func GetCards(cards []string) []Card {
	var ret []Card
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

/*
GetCard - Fetch a card model for a UUID

Parameters:
uuid (string) - The UUID you want a card model for

Returns
result (card.Card) - The card that was found
errors.ErrInvalidUUID - If the UUID is not valid
errors.ErrNoCard - If the card is not found
*/
func GetCard(uuid string) (Card, error) {
	var result Card

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

/*
IndexCards - Return all cards from the database

Parameters:
limit (int64) - Limit the ammount of cards you want returned

Returns:
result (slice[card.Card]) - The results of the operation
errors.ErrNoCards - If the database has no cards
*/
func IndexCards(limit int64) ([]Card, error) {
	var result []Card

	var database = context.GetDatabase()

	results := database.Index("card", limit, &result)
	if results == nil {
		return result, errors.ErrNoCards
	}

	return result, nil

}
