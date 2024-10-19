package deck

import (
	"go.mongodb.org/mongo-driver/bson"
	"mtgjson/context"
	"mtgjson/errors"
	"mtgjson/models/card"
)

type Deck struct {
	Code        string   `json:"code"`
	Commander   []string `json:"commander"`
	MainBoard   []string `json:"mainBoard"`
	Name        string   `json:"name"`
	ReleaseDate string   `json:"releaseDate"`
	SideBoard   []string `json:"sideBoard"`
	Type        string   `json:"type"`
}

func (d Deck) GetMainboard() []card.CardSet {
	return card.GetCards(d.MainBoard)
}

func (d Deck) GetSideboard() []card.CardSet {
	return card.GetCards(d.SideBoard)
}

func (d Deck) GetCommander() []card.CardSet {
	return card.GetCards(d.Commander)
}

func (d Deck) CardExists(uuid string) bool {
	var mainBoard = d.MainBoard
	var ret = false

	for i := 0; i < len(mainBoard); i++ {
		_uuid := mainBoard[i]

		if uuid == _uuid {
			ret = true
			break
		}
	}

	return ret
}

func (d Deck) ValidateCards() (bool, []string) {
	var allCards []string
	var invalidCards []string

	allCards = append(d.MainBoard, d.Commander...)
	allCards = append(allCards, d.SideBoard...)

	result := true
	for i := range allCards {
		uuid := allCards[i]

		_, err := card.GetCard(uuid)
		if err != nil {
			invalidCards = append(invalidCards, uuid)
			result = false
			// not ending the iteration here to ensure the caller is aware of all unidentifiable cards
		}
	}

	return result, invalidCards
}

func (d *Deck) AddCard(uuid string) error {
	var exists = d.CardExists(uuid)
	if exists {
		return errors.ErrCardAlreadyExist
	}

	d.MainBoard = append(d.MainBoard, uuid)

	return nil
}

func (d *Deck) DeleteCard(uuid string) error {
	var exists = d.CardExists(uuid)
	if !exists {
		return errors.ErrNoCard
	}

	var index int
	for i := range d.MainBoard {
		if d.MainBoard[i] == uuid {
			index = i
			break
		}
	}

	d.MainBoard = append(d.MainBoard[:index], d.MainBoard[index+1:]...)

	return nil
}

func (d Deck) UpdateDeck() error {
	var database = context.GetDatabase()

	results := database.Replace("deck", bson.M{"code": d.Code}, &d)
	if results == nil {
		return errors.ErrDeckUpdateFailed
	}

	return nil
}

func GetDeck(code string) (Deck, error) {
	var result Deck

	var database = context.GetDatabase()

	query := bson.M{"code": code}
	results := database.Find("deck", query, &result)
	if results == nil {
		return result, errors.ErrNoDeck
	}

	return result, nil
}

func GetDecks(limit int64) ([]Deck, error) {
	var result []Deck

	var database = context.GetDatabase()

	results := database.Index("deck", limit, &result)
	if results == nil {
		return result, errors.ErrNoDecks
	}

	return result, nil
}

func NewDeck(deck Deck) error {
	if deck.Name == "" || deck.Code == "" {
		return errors.ErrDeckMissingId
	}

	_, valid := GetDeck(deck.Code)
	if valid != errors.ErrNoDeck {
		return errors.ErrDeckAlreadyExists
	}

	var database = context.GetDatabase()

	database.Insert("deck", &deck)

	return nil
}

func (d *Deck) DeleteDeck() any {
	var database = context.GetDatabase()

	query := bson.M{"code": d.Code}
	result := database.Delete("deck", query)
	if result == nil {
		return errors.ErrNoDeck
	}

	if result.DeletedCount != 1 {
		return errors.ErrDeckDeleteFailed
	}

	return result
}
