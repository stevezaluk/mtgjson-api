package deck

import (
	"go.mongodb.org/mongo-driver/bson"
	"mtgjson/context"
	"mtgjson/errors"
	"mtgjson/models/card"
	"mtgjson/server"
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
	return card.IterCards(d.MainBoard)
}

func (d Deck) GetSideboard() []card.CardSet {
	return card.IterCards(d.SideBoard)
}

func (d Deck) GetCommander() []card.CardSet {
	return card.IterCards(d.Commander)
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

func (d Deck) UpdateDeck() error {
	var database = context.ServerContext.Value("database").(server.Database)

	results := database.Replace("deck", bson.M{"code": d.Code}, &d)
	if results == nil {
		return errors.ErrDeckUpdateFailed
	}

	return nil
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

func GetDeck(code string) (Deck, error) {
	var result Deck

	var database = context.ServerContext.Value("database").(server.Database) // create function for this

	query := bson.M{"code": code}
	results := database.Find("deck", query, &result)
	if results == nil {
		return result, errors.ErrNoDeck
	}

	return result, nil
}

func GetDecks(limit int64) ([]Deck, error) {
	var result []Deck

	var database = context.ServerContext.Value("database").(server.Database)

	results := database.Index("deck", limit, &result)
	if results == nil {
		return result, errors.ErrNoDecks
	}

	return result, nil
}
