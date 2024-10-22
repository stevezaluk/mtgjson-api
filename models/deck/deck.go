package deck

import (
	"mtgjson/context"
	"mtgjson/errors"
	"mtgjson/models/card"
	"slices"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	MAINBOARD = "mainBoard"
	SIDEBOARD = "sideBoard"
	COMMANDER = "commanderBoard"
)

/*
Deck - Represents a MTGJSON deck

Code (string) - A 3 or 4 digit code as an identifier for the deck
Commander (slice[string]) - A list of UUID's that represents the commander for the deck
Mainboard (slice[string]) - A list of UUID's that represents the main board for the deck
Name (string) - The name of the deck
ReleaseDate (string) - The release date of the deck
Sideboard (slice[string]) - A list of UUID's that represents the side board for the deck
Type (string) - The deck type
*/
type Deck struct {
	Code        string   `json:"code"`
	Commander   []string `json:"commander"`
	MainBoard   []string `json:"mainBoard"`
	Name        string   `json:"name"`
	ReleaseDate string   `json:"releaseDate"`
	SideBoard   []string `json:"sideBoard"`
	Type        string   `json:"type"`
}

/*
FetchMainboard - Iterate through the UUID's in the main board and return card models

Parameters:
None

Returns
slice[card.Card] - The results
*/
func (d Deck) FetchMainboard() []card.Card {
	return card.GetCards(d.MainBoard)
}

/*
FetchSideboard - Iterate through the UUID's in the side board and return card models

Parameters:
None

Returns
slice[card.Card] - The results
*/
func (d Deck) FetchSideboard() []card.Card {
	return card.GetCards(d.SideBoard)
}

/*
FetchCommander - Iterate through the UUID's in the commander board and return card models

Parameters:
None

Returns
slice[card.Card] - The results
*/
func (d Deck) FetchCommander() []card.Card {
	return card.GetCards(d.Commander)
}

/*
GetBoard - Returns a pointer to the slice that represents the requested board

Parameters:
board (string) - The board you want a pointer too

Returns
*slice[string] - The board the caller requested
*/
func (d *Deck) GetBoard(board string) *[]string {
	if board == MAINBOARD {
		return &d.MainBoard
	} else if board == SIDEBOARD {
		return &d.SideBoard
	} else if board == COMMANDER {
		return &d.Commander
	}

	return nil
}

/*
CardExists - Ensure that a card exists on a specific board using a UUID

Parameters:
uuid (string) - The uuid to check
board (string) - The board you want to check. Can be either: mainBoard, sideBoard, commanderBoard
*/
func (d Deck) CardExists(uuid string, board string) (bool, error) {
	sourceBoard := d.GetBoard(board)
	if sourceBoard == nil {
		return false, errors.ErrBoardNotExist
	}
	var ret = false

	for _, val := range *sourceBoard {
		if val == uuid {
			ret = true
			break
		}
	}
	return ret, nil
}

/*
AllCards - Combine all boards into a list of UUID's

Parameters: None

Returns:
allCard ([]string) - A list of all UUID's in the deck
*/
func (d Deck) AllCards() []string {
	var allCards []string

	allCards = append(d.MainBoard, d.SideBoard...)
	allCards = append(allCards, d.Commander...)

	return allCards
}

/*
AddCards - Add a list of cards to a specific board within a deck.
Card validation is not performed here as it is performed before this operation

Parameters:
uuids (slice[string]) - A list of UUID's you want to add to your deck
board (string) - The board you want to add to

Returns
errors.ErrBoardNotExist - If the board does not exist
*/
func (d *Deck) AddCards(uuids []string, board string) error {
	sourceBoard := d.GetBoard(board)
	if sourceBoard == nil {
		return errors.ErrBoardNotExist
	}

	*sourceBoard = append(*sourceBoard, uuids...)

	return nil
}

/*
DeleteCards - Delete a list of cards to a specific board within a deck.
Card validation is not performed here as it is performed before this operation

Parameters:
uuids (slice[string]) - A list of UUID's you want to remove from your deck
board (string) - The board you want to add to

Returns
errors.ErrBoardNotExist - If the board does not exist
*/
func (d *Deck) DeleteCards(uuids []string, board string) error {
	sourceBoard := d.GetBoard(board)
	if sourceBoard == nil {
		return errors.ErrBoardNotExist
	}

	for _, uuid := range uuids {
		for i, val := range *sourceBoard {
			if uuid == val {
				*sourceBoard = slices.Delete(*sourceBoard, i, i+1)
			}
		}
	}

	return nil
}

/*
UpdateDeck - Replace the deck in the database

Parameters:
None

Returns:
error.ErrDeckUpdateFailed - If database.Replace returns an error
*/
func (d Deck) UpdateDeck() error {
	var database = context.GetDatabase()

	results := database.Replace("deck", bson.M{"code": d.Code}, &d)
	if results == nil {
		return errors.ErrDeckUpdateFailed
	}

	return nil
}

/*
DeleteDeck - Delete the deck from the database

Parameters:
None

Returns:
errors.ErrNoDeck - If the deck does not exist
errors.ErrDeckDeleteFailed - If the mongo results structure doesn't show any deleted results
*/
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

/*
GetDeck - Fetch a deck model and from a deck code

Parameters:
code (string) - The deck code

Returns
Deck (deck.Deck) - A deck model
errors.ErrNoDeck - If the deck does not exist
*/
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

/*
GetDecks - Fetch all decks available in the database

Parameters:
limit (int64) - Limit the ammount of results you want

Returns:
result (slice[deck.Deck]) - The results
errors.ErrNoDecks - If no decks exist in the database
*/
func GetDecks(limit int64) ([]Deck, error) {
	var result []Deck

	var database = context.GetDatabase()

	results := database.Index("deck", limit, &result)
	if results == nil {
		return result, errors.ErrNoDecks
	}

	return result, nil
}

/*
NewDeck - Create a new deck from a deck model

Parameters:
errors.ErrDeskMissingId - If the deck passed in the parameter does not have a valid name or code
errors.ErrDeckAlreadyExists - If the deck already exists under the same code
*/
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
