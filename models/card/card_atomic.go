package models

import (
	"mtgjson/models/card/meta"
)

type CardAtomic struct {
	AsciiName               string                  `json:"asciiName"`
	AttractionLights        []string                `json:"attractionLights"`
	ColorIdentity           []string                `json:"colorIdentity"`
	ColorIndicator          []string                `json:"colorIndicator"`
	Colors                  []string                `json:"colors"`
	ConvertedManaCost       int                     `json:"convertedManaCost"`
	Defense                 string                  `json:"defense"`
	EDHRecRank              int                     `json:"edhrecRank"`
	EDHRecSaltiness         int                     `json:"edhrecSaltiness"`
	FaceConvertedManaCost   int                     `json:"faceConvertedManaCost"`
	FaceManaValue           int                     `json:"faceManaValue"`
	FaceName                string                  `json:"faceName"`
	FirstPrinting           string                  `json:"firstPrinting"`
	ForeignData             models.ForeignData      `json:"foreignData"`
	Hand                    string                  `json:"hand"`
	HasAlternativeDeckLimit bool                    `json:"hasAlternativeDeckLimit"`
	Identifiers             models.CardIdentifiers  `json:"identifiers"`
	IsFunny                 bool                    `json:"isFunny"`
	IsReserved              bool                    `json:"isReserved"`
	Keywords                []string                `json:"keywords"`
	Layout                  string                  `json:"layout"`
	LeadershipSkills        models.LeadershipSkills `json:"leadershipSkills"`
	Legalities              models.CardLegalities   `json:"legalities"`
	Life                    string                  `json:"life"`
	Loyalty                 string                  `json:"loyalty"`
	ManaCost                string                  `json:"manaCost"`
	ManaValue               int                     `json:"manaValue"`
	Name                    string                  `json:"name"`
	Power                   string                  `json:"power"`
	Printings               []string                `json:"printings"`
	PurchaseUrls            models.PurchaseUrls     `json:"purchaseUrls"`
	RelatedCards            models.RelatedCards     `json:"relatedCards"`
	Rulings                 models.CardRulings      `json:"rulings"`
	Side                    string                  `json:"side"`
	Subsets                 []string                `json:"subsets"`
	Subtypes                []string                `json:"subtypes"`
	Supertypes              []string                `json:"supertypes"`
	Text                    string                  `json:"text"`
	Toughness               string                  `json:"toughness"`
	Type                    string                  `json:"type"`
	Types                   []string                `json:"types"`
}
