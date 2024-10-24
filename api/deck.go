package api

import (
	deck_model "github.com/stevezaluk/mtgjson-models/deck"
	"github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/card"
	"github.com/stevezaluk/mtgjson-sdk/deck"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
DeckGET - All logic needed for fetching card metadata

Parameters:
c (gin.Context) - The request context

Returns:
Nothing
*/
func DeckGET(c *gin.Context) {
	code := c.Query("deckCode")
	if code == "" {
		limit := limitToInt64(c.DefaultQuery("limit", "100"))
		results, err := deck.GetDecks(limit)
		if err == errors.ErrNoDecks {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusFound, results)
		return
	}

	results, err := deck.GetDeck(code)
	if err == errors.ErrNoDeck {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusFound, results)
}

/*
DeckPOST - All logic needed for creating new card metadata

Parameters:
c (gin.Context) - The request context

Returns:
Nothing
*/
func DeckPOST(c *gin.Context) {
	var new deck_model.Deck

	if c.Bind(&new) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to bind response to object. Object structure may be incorrect"})
		return
	}

	if new.Name == "" || new.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Deck Code or Name must not be empty when creating a deck"})
		return
	}

	var valid, invalidCards, noExistCards = card.ValidateCards(new.AllCards())
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create deck. Some cards do not exist in the database or are invalid", "invalid": invalidCards, "noExist": noExistCards})
		return
	}

	var err = deck.NewDeck(new)
	if err == errors.ErrDeckAlreadyExists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Deck already exists under this deck code", "deckCode": new.Code})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Successfully created new deck", "deckCode": new.Code})
}

/*
DeckDELETE - All logic needed for deleting card metadata

Parameters:
c (gin.Context) - The request context

Returns:
Nothing
*/
func DeckDELETE(c *gin.Context) {
	code := c.Query("deckCode")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to perform a DELETE operation"})
		return
	}

	_deck, err := deck.GetDeck(code)
	if err == errors.ErrNoDeck {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	result := deck.DeleteDeck(_deck.Code)
	if result == errors.ErrDeckDeleteFailed {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Successfully deleted deck", "deckCode": _deck.Code})
}

/*
DeckContentGET - All logic needed for fetching the contents of a deck

Parameters:
c (gin.Context) - The request context

Returns:
Nothing
*/
func DeckContentGET(c *gin.Context) {
	code := c.Query("deckCode")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to fetch a deck's contents"})
		return
	}

	_deck, err := deck.GetDeck(code)
	if err == errors.ErrNoDeck {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	var mainBoard = deck.FetchMainboard(_deck)
	var sideBoard = deck.FetchSideboard(_deck)
	var commander = deck.FetchCommander(_deck)

	var resp = gin.H{"mainBoard": mainBoard, "sideBoard": sideBoard, "commander": commander}

	c.JSON(http.StatusFound, resp)
}

/*
DeckContentPOST - All logic needed for updating deck contents

Parameters:
c (gin.Context) - The request context

Returns:
Nothing
*/
func DeckContentPOST(c *gin.Context) {
	code := c.Query("deckCode")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to perform a POST operation"})
		return
	}

	_deck, err := deck.GetDeck(code)
	if err == errors.ErrNoDeck {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	var updates deck_model.DeckUpdate
	c.BindJSON(&updates)

	valid, invalidCards, noExistCards := card.ValidateCards(updates.AllCards())
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update deck. Some cards do not exist in the database or are invalid", "invalid": invalidCards, "noExist": noExistCards})
		return
	}

	_deck.AddCards(updates.MainBoard, deck_model.MAINBOARD)
	_deck.AddCards(updates.SideBoard, deck_model.SIDEBOARD)
	_deck.AddCards(updates.Commander, deck_model.COMMANDER)

	err = deck.UpdateDeck(_deck)
	if err == errors.ErrDeckUpdateFailed {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Successfully updated deck", "deckCode": code, "count": updates.Count()})
}

/*
DeckContentDELETE - All logic needed for removing cards from the deck

Parameters:
c (gin.Context) - The request context

Returns:
Nothing
*/
func DeckContentDELETE(c *gin.Context) {
	code := c.Query("deckCode")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to perform a DELETE operation"})
		return
	}

	_deck, err := deck.GetDeck(code)
	if err == errors.ErrNoDeck {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	var updates deck_model.DeckUpdate
	c.BindJSON(&updates)

	_deck.DeleteCards(updates.MainBoard, deck_model.MAINBOARD)
	_deck.DeleteCards(updates.SideBoard, deck_model.SIDEBOARD)
	_deck.DeleteCards(updates.Commander, deck_model.COMMANDER)

	err = deck.UpdateDeck(_deck)
	if err == errors.ErrDeckUpdateFailed {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Successfully removed cards from deck", "deckCode": code, "count": updates.Count()})
}
